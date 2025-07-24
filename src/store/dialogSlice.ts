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

export const dialogActions = {
  openDialog: (payload: { type: string; payload?: any }) => ({
    type: 'dialog/openDialog',
    payload
  }),
  openAlertDialog: (payload: { type: string; payload: any }) => ({
    type: 'dialog/openAlertDialog',
    payload
  }),
  openVerifyWhatsAppDialog: (payload: { type: string; payload: { phoneNumber: string; token: string } }) => ({
    type: 'dialog/openVerifyWhatsAppDialog',
    payload
  }),
  closeDialog: () => ({
    type: 'dialog/closeDialog'
  })
};

const dialogReducer = (state = initialState, action: any): DialogState => {
  switch (action.type) {
    case 'dialog/openDialog':
      return {
        isOpen: true,
        type: action.payload.type,
        payload: action.payload.payload || null
      };
    case 'dialog/openAlertDialog':
      return {
        isOpen: true,
        type: action.payload.type,
        payload: action.payload.payload
      };
    case 'dialog/openVerifyWhatsAppDialog':
      return {
        isOpen: true,
        type: action.payload.type,
        payload: action.payload.payload
      };
    case 'dialog/closeDialog':
      return {
        isOpen: false,
        type: null,
        payload: null
      };
    default:
      return state;
  }
};

export default dialogReducer;
