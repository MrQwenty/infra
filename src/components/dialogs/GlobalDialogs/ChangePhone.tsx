import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';

import { dialogActions } from '../../../store/dialogSlice';
import { userActions } from '../../../store/userSlice';
import { getUserReq, changePhoneReq } from '../../../api/userAPI';

interface ChangePhoneProps {
  onClose: () => void;
}

const ChangePhone: React.FC<ChangePhoneProps> = ({ onClose }) => {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const [formData, setFormData] = useState({
    newPhone: ''
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await changePhoneReq(formData.newPhone);
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
          }
        }))
      } else {
        setError(response.error || 'Failed to change phone number');
      }
    } catch (e: any) {
      setError(e.message || 'An error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="dialog-overlay">
      <div className="dialog-content">
        <h2>{t('changePhone.title')}</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="newPhone">{t('changePhone.newPhoneLabel')}</label>
            <input
              type="tel"
              id="newPhone"
              name="newPhone"
              value={formData.newPhone}
              onChange={handleInputChange}
              placeholder={t('changePhone.newPhonePlaceholder')}
              required
            />
          </div>
          {error && <div className="error-message">{error}</div>}
          <div className="dialog-actions">
            <button type="button" onClick={onClose} disabled={isLoading}>
              {t('changePhone.cancelBtn')}
            </button>
            <button type="submit" disabled={isLoading}>
              {isLoading ? t('changePhone.changing') : t('changePhone.changeBtn')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ChangePhone;
