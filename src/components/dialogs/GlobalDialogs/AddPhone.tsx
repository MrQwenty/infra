import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux';

import { dialogActions } from '../../../store/dialogSlice';
import { userActions } from '../../../store/userSlice';
import { getUserReq } from '../../../api/userAPI';

interface AddPhoneProps {
  onClose: () => void;
}

const AddPhone: React.FC<AddPhoneProps> = ({ onClose }) => {
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
      const response = await fetch('/v1/user/phone/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          phoneNumber: formData.newPhone
        })
      });

      if (!response.ok) {
        throw new Error('Failed to add phone number');
      }

      const data = await response.json();

      if (data.success) {
        const userData = (await getUserReq()).data;
        dispatch(userActions.setUser(userData));
      }
        dispatch(dialogActions.openVerifyWhatsAppDialog({
          type: 'verifyWhatsApp',
          payload: {
            phoneNumber: formData.newPhone,
            token: data.verificationToken || '',
          }
        }))

    } catch (e: any) {
      setError(e.message || 'An error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="dialog-overlay">
      <div className="dialog-content">
        <h2>{t('addPhone.title')}</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="newPhone">{t('addPhone.newPhoneLabel')}</label>
            <input
              type="tel"
              id="newPhone"
              name="newPhone"
              value={formData.newPhone}
              onChange={handleInputChange}
              placeholder={t('addPhone.newPhonePlaceholder')}
              required
            />
          </div>
          {error && <div className="error-message">{error}</div>}
          <div className="dialog-actions">
            <button type="button" onClick={onClose} disabled={isLoading}>
              {t('addPhone.cancelBtn')}
            </button>
            <button type="submit" disabled={isLoading}>
              {isLoading ? t('addPhone.adding') : t('addPhone.addBtn')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddPhone;
