import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';

import { dialogActions } from '../../../store/dialogSlice';
import { verifyWhatsAppReq, resendWhatsAppCodeReq, cancelWhatsAppVerificationReq } from '../../../api/userAPI';

interface VerifyWhatsAppProps {
  phoneNumber: string;
  token: string;
  verificationMethod?: 'whatsapp' | 'sms';
  onClose: () => void;
}

const VerifyWhatsApp: React.FC<VerifyWhatsAppProps> = ({ 
  phoneNumber, 
  token, 
  verificationMethod = 'whatsapp', 
  onClose 
}) => {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const [verificationCode, setVerificationCode] = useState('');
  const [isVerifying, setIsVerifying] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [error, setError] = useState('');
  const [attemptsRemaining, setAttemptsRemaining] = useState<number | null>(null);
  const [timeRemaining, setTimeRemaining] = useState<number>(600); // 10 minutes in seconds

  // Countdown timer effect
  React.useEffect(() => {
    const timer = setInterval(() => {
      setTimeRemaining(prev => {
        if (prev <= 1) {
          clearInterval(timer);
          setError(t('dialogs.verifyWhatsApp.errors.expired'));
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [t]);

  const formatTime = (seconds: number): string => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const handleVerify = async () => {
    if (!verificationCode.trim()) {
      setError(t('dialogs.verifyWhatsApp.errors.codeRequired'));
      return;
    }

    // Validate code format (6 digits)
    if (!/^\d{6}$/.test(verificationCode.trim())) {
      setError(t('dialogs.verifyWhatsApp.errors.invalidFormat'));
      return;
    }

    setIsVerifying(true);
    setError('');

    try {
      const response = await verifyWhatsAppReq(token, verificationCode);
      
      if (response.success) {
        // Show success message briefly before closing
        setError(t('dialogs.verifyWhatsApp.success'));
        setTimeout(() => {
          dispatch(dialogActions.closeDialog());
        }, 1000);
      } else {
        if (response.attemptsRemaining !== undefined) {
          setAttemptsRemaining(response.attemptsRemaining);
          if (response.attemptsRemaining === 0) {
            setError(t('dialogs.verifyWhatsApp.errors.maxAttemptsReached'));
            setTimeout(() => {
              handleCancel();
            }, 3000);
          } else {
            setError(t('dialogs.verifyWhatsApp.errors.verificationFailed', { 
              attempts: response.attemptsRemaining 
            }));
          }
        } else {
          setError(t('dialogs.verifyWhatsApp.errors.verificationFailed'));
        }
      }
    } catch (err: any) {
      if (err.message.includes('expired')) {
        setError(t('dialogs.verifyWhatsApp.errors.expired'));
      } else if (err.message.includes('attempts')) {
        setError(t('dialogs.verifyWhatsApp.errors.maxAttemptsReached'));
        setTimeout(() => {
          handleCancel();
        }, 3000);
      } else {
        setError(t('dialogs.verifyWhatsApp.errors.verificationFailed'));
      }
    } finally {
      setIsVerifying(false);
      setVerificationCode(''); // Clear the input for security
    }
  };

  const handleCancel = async () => {
    try {
      await cancelWhatsAppVerificationReq(token);
    } catch (err) {
      console.error('Failed to cancel verification:', err);
    } finally {
      dispatch(dialogActions.closeDialog());
    }
  };

  const handleResend = async () => {
    setIsResending(true);
    setError('');

    try {
      const response = await resendWhatsAppCodeReq(token);
      if (response.success) {
        setTimeRemaining(600); // Reset timer to 10 minutes
        setAttemptsRemaining(response.attemptsRemaining);
        setVerificationCode(''); // Clear current input
        // Don't show success message to avoid confusion
      } else {
        setError(t('dialogs.verifyWhatsApp.errors.resendFailed'));
      }
    } catch (err: any) {
      if (err.message.includes('expired')) {
        setError(t('dialogs.verifyWhatsApp.errors.expired'));
      } else {
        setError(t('dialogs.verifyWhatsApp.errors.resendFailed'));
      }
    } finally {
      setIsResending(false);
    }
  };

  // Auto-format verification code input
  const handleCodeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
    setVerificationCode(value);
    setError(''); // Clear error when user starts typing
  };

  // Auto-submit when 6 digits are entered
  React.useEffect(() => {
    if (verificationCode.length === 6 && !isVerifying) {
      handleVerify();
    }
  }, [verificationCode]);

  const isExpired = timeRemaining === 0;
  const methodName = verificationMethod === 'whatsapp' ? 'WhatsApp' : 'SMS';

  return (
    <div className="dialog-overlay">
      <div className="dialog-content verification-dialog">
        <h2>{t('dialogs.verifyWhatsApp.title')}</h2>
        <p className="verification-description">
          {t('dialogs.verifyWhatsApp.description', { phoneNumber, method: methodName })}
        </p>
        
        <div className="verification-timer">
          <span className={`timer ${timeRemaining < 60 ? 'warning' : ''}`}>
            {t('dialogs.verifyWhatsApp.timeRemaining')}: {formatTime(timeRemaining)}
          </span>
        </div>
        
        <div className="form-group">
          <label htmlFor="verificationCode">{t('dialogs.verifyWhatsApp.codeInputLabel')}</label>
          <input
            type="text"
            id="verificationCode"
            value={verificationCode}
            onChange={handleCodeChange}
            placeholder={t('dialogs.verifyWhatsApp.codeInputPlaceholder')}
            maxLength={6}
            pattern="[0-9]{6}"
            className="verification-code-input"
            disabled={isVerifying || isExpired}
            autoComplete="one-time-code"
            inputMode="numeric"
          />
          <div className="input-hint">
            {verificationCode.length}/6 {t('dialogs.verifyWhatsApp.digitsEntered')}
          </div>
        </div>

        {attemptsRemaining !== null && attemptsRemaining > 0 && (
          <div className="attempts-remaining">
            {t('dialogs.verifyWhatsApp.attemptsRemaining', { count: attemptsRemaining })}
          </div>
        )}

        {error && (
          <div className={`error-message ${error.includes('success') ? 'success-message' : ''}`}>
            {error}
          </div>
        )}

        <div className="dialog-actions">
          <button type="button" onClick={handleCancel} disabled={isVerifying || isResending}>
            {t('dialogs.verifyWhatsApp.cancelBtn')}
          </button>
          
          <button 
            type="button" 
            onClick={handleResend} 
            disabled={isResending || isVerifying || isExpired}
            className="secondary"
          >
            {isResending ? t('dialogs.verifyWhatsApp.resending') : t('dialogs.verifyWhatsApp.resendBtn')}
          </button>
          
          <button 
            type="button" 
            onClick={handleVerify} 
            disabled={isVerifying || !verificationCode.trim() || verificationCode.length !== 6 || isExpired}
            className="primary"
          >
            {isVerifying ? t('dialogs.verifyWhatsApp.verifying') : t('dialogs.verifyWhatsApp.verifyBtn')}
          </button>
        </div>
      </div>
    </div>
  );
};

export default VerifyWhatsApp;