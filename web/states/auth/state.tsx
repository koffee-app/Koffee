export interface IUser {
  id: number;
  email: string;
}

export interface IUserError {
  email?: string;
  password?: string;
}

export interface IAuth {
  authenticated: boolean;
  user?: IUser;
}
