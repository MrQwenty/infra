import { apiBase } from '../constants';
import { WhatsAppVerificationRequest, WhatsAppVerificationResponse, VerifyCodeRequest, VerifyCodeResponse } from '../types/verification';

export interface User {
  id: string;
  email: string;
  phoneNumber?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  error?: string;
  user?: User;
}

export const getUserReq = async (): Promise<ApiResponse<User>> => {
  const token = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/profile`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch user data');
  }

  return response.json();
};

export const addPhoneReq = async (phoneNumber: string, verificationMethod: 'whatsapp' | 'sms' = 'whatsapp'): Promise<ApiResponse<{ verificationToken: string }>> => {
  const token = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/phone/add`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token,
    },
    body: JSON.stringify({ phoneNumber, verificationMethod }),
  });

  if (!response.ok) {
    throw new Error('Failed to add phone number');
  }

  return response.json();
};

export const changePhoneReq = async (newPhoneNumber: string, verificationMethod: 'whatsapp' | 'sms' = 'whatsapp'): Promise<ApiResponse<{ verificationToken: string }>> => {
  const token = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/phone/change`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token,
    },
    body: JSON.stringify({ newPhoneNumber, verificationMethod }),
  });

  if (!response.ok) {
    throw new Error('Failed to change phone number');
  }

  return response.json();
};

export const initiateWhatsAppVerificationReq = async (request: WhatsAppVerificationRequest): Promise<WhatsAppVerificationResponse> => {
  const token = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/whatsapp-verification/initiate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token,
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error('Failed to initiate WhatsApp verification');
  }

  return response.json();
};

export const verifyWhatsAppReq = async (token: string, code: string): Promise<ApiResponse<{}>> => {
  const authToken = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/whatsapp-verification/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': authToken,
    },
    body: JSON.stringify({ token, code }),
  });

  if (!response.ok) {
    throw new Error('Failed to verify WhatsApp code');
  }

  return response.json();
};

export const resendWhatsAppCodeReq = async (token: string): Promise<WhatsAppVerificationResponse> => {
  const authToken = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/whatsapp-verification/resend`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': authToken,
    },
    body: JSON.stringify({ token }),
  });

  if (!response.ok) {
    throw new Error('Failed to resend WhatsApp code');
  }

  return response.json();
};

export const cancelWhatsAppVerificationReq = async (token: string): Promise<ApiResponse<{}>> => {
  const authToken = localStorage.getItem('authToken') || '';
  const response = await fetch(`${apiBase}/user/whatsapp-verification/cancel`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': authToken,
    },
    body: JSON.stringify({ token }),
  });

  if (!response.ok) {
    throw new Error('Failed to cancel WhatsApp verification');
  }

  return response.json();
};