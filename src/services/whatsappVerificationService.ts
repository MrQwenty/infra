import { VerificationAttempt, WhatsAppVerificationRequest, WhatsAppVerificationResponse, VerifyCodeRequest, VerifyCodeResponse } from '../types/verification';

class WhatsAppVerificationService {
  private readonly MAX_ATTEMPTS = 3;
  private readonly MAX_RETRIES = 3;
  private readonly CODE_EXPIRY_MINUTES = 10;
  private readonly RETRY_DELAYS = [30000, 60000, 120000]; // 30s, 1m, 2m
  
  private verificationAttempts: Map<string, VerificationAttempt> = new Map();
  private retryTimeouts: Map<string, NodeJS.Timeout> = new Map();

  /**
   * Generate a 6-digit verification code
   */
  private generateVerificationCode(): string {
    return Math.floor(100000 + Math.random() * 900000).toString();
  }

  /**
   * Generate a unique verification token
   */
  private generateVerificationToken(): string {
    const timestamp = Date.now();
    const random = Math.random().toString(36).substring(2);
    return `whatsapp_${timestamp}_${random}`;
  }

  /**
   * Send WhatsApp verification code with retry mechanism
   */
  private async sendWhatsAppMessage(phoneNumber: string, code: string, retryCount: number = 0): Promise<boolean> {
    try {
      // Simulate WhatsApp API call
      const response = await fetch('/api/whatsapp/send-verification', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${process.env.WHATSAPP_API_TOKEN}`,
        },
        body: JSON.stringify({
          to: phoneNumber,
          template: 'phone_verification',
          components: [{
            type: 'body',
            parameters: [{
              type: 'text',
              text: code
            }]
          }]
        })
      });

      if (response.ok) {
        console.log(`WhatsApp verification sent successfully to ${phoneNumber}`);
        return true;
      } else {
        throw new Error(`WhatsApp API error: ${response.status}`);
      }
    } catch (error) {
      console.error(`Failed to send WhatsApp verification (attempt ${retryCount + 1}):`, error);
      
      if (retryCount < this.MAX_RETRIES) {
        const delay = this.RETRY_DELAYS[retryCount] || this.RETRY_DELAYS[this.RETRY_DELAYS.length - 1];
        console.log(`Retrying WhatsApp send in ${delay}ms...`);
        
        return new Promise((resolve) => {
          const timeoutId = setTimeout(async () => {
            const result = await this.sendWhatsAppMessage(phoneNumber, code, retryCount + 1);
            resolve(result);
          }, delay);
          
          // Store timeout for potential cleanup
          this.retryTimeouts.set(`${phoneNumber}_${retryCount}`, timeoutId);
        });
      }
      
      return false;
    }
  }

  /**
   * Initiate phone number verification
   */
  async initiateVerification(request: WhatsAppVerificationRequest): Promise<WhatsAppVerificationResponse> {
    const { phoneNumber, verificationMethod } = request;
    
    // Validate phone number format
    const phoneRegex = /^\+[1-9]\d{1,14}$/;
    if (!phoneRegex.test(phoneNumber)) {
      throw new Error('Invalid phone number format');
    }

    // Check if there's an existing active verification
    const existingToken = Array.from(this.verificationAttempts.entries())
      .find(([_, attempt]) => attempt.phoneNumber === phoneNumber && attempt.status === 'pending')?.[0];
    
    if (existingToken) {
      const existing = this.verificationAttempts.get(existingToken)!;
      if (existing.expiresAt > new Date()) {
        return {
          success: true,
          verificationToken: existingToken,
          message: 'Verification code already sent',
          expiresAt: existing.expiresAt,
          attemptsRemaining: existing.maxAttempts - existing.attempts
        };
      } else {
        // Clean up expired verification
        this.verificationAttempts.delete(existingToken);
      }
    }

    const code = this.generateVerificationCode();
    const token = this.generateVerificationToken();
    const expiresAt = new Date(Date.now() + this.CODE_EXPIRY_MINUTES * 60 * 1000);

    const verificationAttempt: VerificationAttempt = {
      id: token,
      phoneNumber,
      code,
      token,
      attempts: 0,
      maxAttempts: this.MAX_ATTEMPTS,
      createdAt: new Date(),
      expiresAt,
      status: 'pending',
      retryCount: 0,
      maxRetries: this.MAX_RETRIES
    };

    this.verificationAttempts.set(token, verificationAttempt);

    // Send verification code via WhatsApp
    const sent = await this.sendWhatsAppMessage(phoneNumber, code);
    
    if (!sent) {
      // If all retries failed, mark as failed and clean up
      verificationAttempt.status = 'failed';
      this.verificationAttempts.delete(token);
      throw new Error('Failed to send verification code after multiple attempts');
    }

    // Set up automatic cleanup after expiry
    setTimeout(() => {
      const attempt = this.verificationAttempts.get(token);
      if (attempt && attempt.status === 'pending') {
        attempt.status = 'expired';
        this.verificationAttempts.delete(token);
        console.log(`Verification token ${token} expired and cleaned up`);
      }
    }, this.CODE_EXPIRY_MINUTES * 60 * 1000);

    return {
      success: true,
      verificationToken: token,
      message: 'Verification code sent via WhatsApp',
      expiresAt,
      attemptsRemaining: this.MAX_ATTEMPTS
    };
  }

  /**
   * Verify the entered code
   */
  async verifyCode(request: VerifyCodeRequest): Promise<VerifyCodeResponse> {
    const { token, code } = request;
    
    const verificationAttempt = this.verificationAttempts.get(token);
    
    if (!verificationAttempt) {
      return {
        success: false,
        message: 'Invalid or expired verification token',
        verified: false
      };
    }

    // Check if verification has expired
    if (verificationAttempt.expiresAt < new Date()) {
      verificationAttempt.status = 'expired';
      this.verificationAttempts.delete(token);
      return {
        success: false,
        message: 'Verification code has expired',
        verified: false
      };
    }

    // Check if max attempts reached
    if (verificationAttempt.attempts >= verificationAttempt.maxAttempts) {
      verificationAttempt.status = 'failed';
      this.verificationAttempts.delete(token);
      return {
        success: false,
        message: 'Maximum verification attempts exceeded',
        verified: false
      };
    }

    // Increment attempt counter
    verificationAttempt.attempts++;

    // Verify the code
    if (verificationAttempt.code === code.trim()) {
      verificationAttempt.status = 'verified';
      
      // Clean up successful verification after a short delay
      setTimeout(() => {
        this.verificationAttempts.delete(token);
      }, 5000);

      return {
        success: true,
        message: 'Phone number verified successfully',
        verified: true
      };
    } else {
      const attemptsRemaining = verificationAttempt.maxAttempts - verificationAttempt.attempts;
      
      if (attemptsRemaining === 0) {
        verificationAttempt.status = 'failed';
        this.verificationAttempts.delete(token);
        return {
          success: false,
          message: 'Invalid verification code. Maximum attempts exceeded.',
          verified: false
        };
      }

      return {
        success: false,
        message: `Invalid verification code. ${attemptsRemaining} attempts remaining.`,
        verified: false,
        attemptsRemaining
      };
    }
  }

  /**
   * Resend verification code
   */
  async resendVerificationCode(token: string): Promise<WhatsAppVerificationResponse> {
    const verificationAttempt = this.verificationAttempts.get(token);
    
    if (!verificationAttempt) {
      throw new Error('Invalid or expired verification token');
    }

    if (verificationAttempt.expiresAt < new Date()) {
      verificationAttempt.status = 'expired';
      this.verificationAttempts.delete(token);
      throw new Error('Verification session has expired');
    }

    // Generate new code and reset attempts
    const newCode = this.generateVerificationCode();
    verificationAttempt.code = newCode;
    verificationAttempt.attempts = 0;
    verificationAttempt.retryCount = 0;

    // Extend expiry time
    verificationAttempt.expiresAt = new Date(Date.now() + this.CODE_EXPIRY_MINUTES * 60 * 1000);

    const sent = await this.sendWhatsAppMessage(verificationAttempt.phoneNumber, newCode);
    
    if (!sent) {
      throw new Error('Failed to resend verification code');
    }

    return {
      success: true,
      verificationToken: token,
      message: 'New verification code sent via WhatsApp',
      expiresAt: verificationAttempt.expiresAt,
      attemptsRemaining: this.MAX_ATTEMPTS
    };
  }

  /**
   * Get verification status
   */
  getVerificationStatus(token: string): VerificationAttempt | null {
    return this.verificationAttempts.get(token) || null;
  }

  /**
   * Cancel verification and clean up
   */
  cancelVerification(token: string): boolean {
    const verificationAttempt = this.verificationAttempts.get(token);
    if (verificationAttempt) {
      verificationAttempt.status = 'failed';
      this.verificationAttempts.delete(token);
      
      // Clean up any pending retry timeouts
      this.retryTimeouts.forEach((timeout, key) => {
        if (key.startsWith(verificationAttempt.phoneNumber)) {
          clearTimeout(timeout);
          this.retryTimeouts.delete(key);
        }
      });
      
      return true;
    }
    return false;
  }

  /**
   * Clean up expired verifications (should be called periodically)
   */
  cleanupExpiredVerifications(): void {
    const now = new Date();
    const expiredTokens: string[] = [];

    this.verificationAttempts.forEach((attempt, token) => {
      if (attempt.expiresAt < now || attempt.status !== 'pending') {
        expiredTokens.push(token);
      }
    });

    expiredTokens.forEach(token => {
      this.verificationAttempts.delete(token);
    });

    if (expiredTokens.length > 0) {
      console.log(`Cleaned up ${expiredTokens.length} expired verification attempts`);
    }
  }
}

export const whatsappVerificationService = new WhatsAppVerificationService();

// Set up periodic cleanup (every 5 minutes)
setInterval(() => {
  whatsappVerificationService.cleanupExpiredVerifications();
}, 5 * 60 * 1000);