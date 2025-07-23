import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';
import { dialogActions } from '../../../store/dialogSlice';

interface VerifyWhatsAppProps {
  phoneNumber: string;
  token: string;
  onClose: () => void;
}

const VerifyWhatsApp: React.FC<VerifyWhatsAppProps> = ({ phoneNumber, token, onClose }) => {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const [verificationCode, setVerificationCode] = useState('');
  const [isVerifying, setIsVerifying] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [error, setError] = useState('');

  const handleVerify = async () => {
    if (!verificationCode.trim()) {
      setError(t('verifyWhatsApp.errors.codeRequired'));
      return;
    }

    setIsVerifying(true);
    setError('');

    try {
      const response = await fetch('/v1/user/whatsapp-verification/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token,
          code: verificationCode
        })
      });

      if (!response.ok) {
        throw new Error('Verification failed');
      }

      const data = await response.json();
      
      if (data.success) {
        dispatch(dialogActions.closeDialog());
      } else {
        setError(t('verifyWhatsApp.errors.verificationFailed'));
      }
    } catch (err) {
      setError(t('verifyWhatsApp.errors.verificationFailed'));
    } finally {
      setIsVerifying(false);
    }
  };

  const handleResend = async () => {
    setIsResending(true);
    setError('');

    try {
      const response = await fetch('/v1/user/whatsapp-verification/send', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          phoneNumber,
          token
        })
      });

      if (!response.ok) {
        throw new Error('Failed to resend code');
      }
    } catch (err) {
      setError(t('verifyWhatsApp.errors.resendFailed'));
    } finally {
      setIsResending(false);
    }
  };

  return (
    <div className="dialog-overlay">
      <div className="dialog-content">
        <h2>{t('verifyWhatsApp.title')}</h2>
        <p>{t('verifyWhatsApp.description', { phoneNumber })}</p>
        
        <div className="form-group">
          <label htmlFor="verificationCode">{t('verifyWhatsApp.codeInputLabel')}</label>
          <input
            type="text"
            id="verificationCode"
            value={verificationCode}
            onChange={(e) => setVerificationCode(e.target.value)}
            placeholder={t('verifyWhatsApp.codeInputPlaceholder')}
            maxLength={6}
            pattern="[0-9]{6}"
          />
        </div>

        {error && <div className="error-message">{error}</div>}

        <div className="dialog-actions">
          <button type="button" onClick={onClose}>
            {t('verifyWhatsApp.cancelBtn')}
          </button>
          <button 
            type="button" 
            onClick={handleResend} 
            disabled={isResending}
            className="secondary"
          >
            {isResending ? t('verifyWhatsApp.resending') : t('verifyWhatsApp.resendBtn')}
          </button>
          <button 
            type="button" 
            onClick={handleVerify} 
            disabled={isVerifying || !verificationCode.trim()}
            className="primary"
          >
            {isVerifying ? t('verifyWhatsApp.verifying') : t('verifyWhatsApp.verifyBtn')}
          </button>
        </div>
      </div>
    </div>
  );
};

export default VerifyWhatsApp;
