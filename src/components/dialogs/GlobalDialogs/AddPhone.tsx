import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';

import { dialogActions } from '../../../store/dialogSlice';
import { userActions } from '../../../store/userSlice';
import { getUserReq, addPhoneReq } from '../../../api/userAPI';

interface AddPhoneProps {
  onClose: () => void;
}

const AddPhone: React.FC<AddPhoneProps> = ({ onClose }) => {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const [formData, setFormData] = useState({
    newPhone: '',
    verificationMethod: 'whatsapp' as 'whatsapp' | 'sms'
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [showMethodSelection, setShowMethodSelection] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleMethodChange = (method: 'whatsapp' | 'sms') => {
    setFormData(prev => ({
      ...prev,
      verificationMethod: method
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validate phone number format
    const phoneRegex = /^\+[1-9]\d{1,14}$/;
    if (!phoneRegex.test(formData.newPhone)) {
      setError(t('dialogs.addPhone.errors.wrongPhoneFormat'));
      return;
    }
    
    setShowMethodSelection(true);
  };

  const handleConfirmVerification = async () => {
    setIsLoading(true);
    setError('');

    try {
      const response = await addPhoneReq(formData.newPhone, formData.verificationMethod);
      if (response.success) {
        if (response.user) {
          dispatch(userActions.setUser(response.user));
        } else {
          const userData = (await getUserReq()).data;
          dispatch(userActions.setUser(userData));
        }
        dispatch(dialogActions.openVerifyWhatsAppDialog({
          type: 'verifyWhatsApp',
          payload: {
            phoneNumber: formData.newPhone,
            token: response.data.verificationToken || '',
            verificationMethod: formData.verificationMethod,
          }
        }));
      } else {
        setError(response.error || 'Failed to add phone number');
      }
    } catch (e: any) {
      setError(e.message || 'An error occurred');
    } finally {
      setIsLoading(false);
      setShowMethodSelection(false);
    }
  };

  if (showMethodSelection) {
    return (
      <div className="dialog-overlay">
        <div className="dialog-content">
          <h2>{t('dialogs.addPhone.warningDialog.title')}</h2>
          <p>{t('dialogs.addPhone.warningDialog.content')}</p>
          
          <div className="form-group">
            <label>{t('dialogs.addPhone.warningDialog.verificationMethodLabel')}</label>
            <div className="radio-group">
              <label className="radio-option">
                <input
                  type="radio"
                  name="verificationMethod"
                  value="sms"
                  checked={formData.verificationMethod === 'sms'}
                  onChange={() => handleMethodChange('sms')}
                />
                {t('dialogs.addPhone.warningDialog.smsOption')}
              </label>
              <label className="radio-option">
                <input
                  type="radio"
                  name="verificationMethod"
                  value="whatsapp"
                  checked={formData.verificationMethod === 'whatsapp'}
                  onChange={() => handleMethodChange('whatsapp')}
                />
                {t('dialogs.addPhone.warningDialog.whatsappOption')}
              </label>
            </div>
          </div>

          {error && <div className="error-message">{error}</div>}
          
          <div className="dialog-actions">
            <button type="button" onClick={() => setShowMethodSelection(false)} disabled={isLoading}>
              {t('dialogs.addPhone.warningDialog.cancelBtn')}
            </button>
            <button type="button" onClick={handleConfirmVerification} disabled={isLoading}>
              {isLoading ? t('dialogs.addPhone.adding') : t('dialogs.addPhone.warningDialog.confirmBtn')}
            </button>
          </div>
        </div>
      </div>
    );
  }
  return (
    <div className="dialog-overlay">
      <div className="dialog-content">
        <h2>{t('dialogs.addPhone.title')}</h2>
        <p>{t('dialogs.addPhone.info')}</p>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="newPhone">{t('dialogs.addPhone.phoneInputLabel')}</label>
            <input
              type="tel"
              id="newPhone"
              name="newPhone"
              value={formData.newPhone}
              onChange={handleInputChange}
              placeholder={t('dialogs.addPhone.phoneInputPlaceholder')}
              required
            />
          </div>
          {error && <div className="error-message">{error}</div>}
          <div className="dialog-actions">
            <button type="button" onClick={onClose} disabled={isLoading}>
              {t('dialogs.addPhone.cancelBtn')}
            </button>
            <button type="submit" disabled={isLoading}>
              {t('dialogs.addPhone.confirmBtn')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddPhone;
