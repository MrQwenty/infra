import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface DialogState {
  isOpen: boolean;
  type: string | null;
  payload: any;
}

const initialState: DialogState = {
  isOpen: false,
  type: null,
  payload: null,
};

const dialogSlice = createSlice({
  name: 'dialog',
  initialState,
  reducers: {
    openDialog: (state, action: PayloadAction<{ type: string; payload?: any }>) => {
      state.isOpen = true;
      state.type = action.payload.type;
      state.payload = action.payload.payload || null;
    },
    openAlertDialog: (state, action: PayloadAction<{ type: string; payload: any }>) => {
      state.isOpen = true;
      state.type = action.payload.type;
      state.payload = action.payload.payload;
    },
    openVerifyWhatsAppDialog: (state, action: PayloadAction<{ type: string; payload: { phoneNumber: string; token: string } }>) => {
      state.isOpen = true;
      state.type = action.payload.type;
      state.payload = action.payload.payload;
    },
    closeDialog: (state) => {
      state.isOpen = false;
      state.type = null;
      state.payload = null;
    },
  },
});

export const dialogActions = dialogSlice.actions;
export default dialogSlice.reducer;
