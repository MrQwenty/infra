export interface VerificationAttempt {
  id: string;
  phoneNumber: string;
  code: string;
  token: string;
  attempts: number;
  maxAttempts: number;
  createdAt: Date;
  expiresAt: Date;
  status: 'pending' | 'verified' | 'expired' | 'failed';
  retryCount: number;
  maxRetries: number;
}

export interface WhatsAppVerificationRequest {
  phoneNumber: string;
  verificationMethod: 'whatsapp' | 'sms';
}

export interface WhatsAppVerificationResponse {
  success: boolean;
  verificationToken: string;
  message: string;
  expiresAt: Date;
  attemptsRemaining: number;
}

export interface VerifyCodeRequest {
  token: string;
  code: string;
}

export interface VerifyCodeResponse {
  success: boolean;
  message: string;
  verified: boolean;
  attemptsRemaining?: number;
}