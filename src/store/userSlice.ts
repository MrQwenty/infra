export interface UserState {
  user: any;
}

const initialState: UserState = {
  user: null,
};

export const userActions = {
  setUser: (user: any) => ({
    type: 'user/setUser',
    payload: user
  })
};

const userReducer = (state = initialState, action: any): UserState => {
  switch (action.type) {
    case 'user/setUser':
      return {
        ...state,
        user: action.payload
      };
    default:
      return state;
  }
};

export default userReducer;
