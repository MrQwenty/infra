import { apiBase } from '../constants';

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
  const response = await fetch(`${apiBase}/user/profile`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch user data');
  }

  return response.json();
};

export const addPhoneReq = async (phoneNumber: string): Promise<ApiResponse<{ verificationToken: string }>> => {
  const response = await fetch(`${apiBase}/user/phone/add`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ phoneNumber }),
  });

  if (!response.ok) {
    throw new Error('Failed to add phone number');
  }

  return response.json();
};

export const changePhoneReq = async (newPhoneNumber: string): Promise<ApiResponse<{ verificationToken: string }>> => {
  const response = await fetch(`${apiBase}/user/phone/change`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ newPhoneNumber }),
  });

  if (!response.ok) {
    throw new Error('Failed to change phone number');
  }

  return response.json();
};

export const verifyWhatsAppReq = async (token: string, code: string): Promise<ApiResponse<{}>> => {
  const response = await fetch(`${apiBase}/user/whatsapp-verification/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ token, code }),
  });

  if (!response.ok) {
    throw new Error('Failed to verify WhatsApp code');
  }

  return response.json();
};

export const resendWhatsAppCodeReq = async (phoneNumber: string, token: string): Promise<ApiResponse<{}>> => {
  const response = await fetch(`${apiBase}/user/whatsapp-verification/send`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ phoneNumber, token }),
  });

  if (!response.ok) {
    throw new Error('Failed to resend WhatsApp code');
  }

  return response.json();
};
