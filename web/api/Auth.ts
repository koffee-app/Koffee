import { IUserForm } from 'states/auth/actions';
import { IUser, IUserError } from 'states/auth/state';
import { whether, Whether } from '../lib/some';

class Auth {
  public static register(form: IUserForm): Whether<IUser, IUserError> {
    return whether({ email: form.email, id: 3 }, null);
  }
}

export default Auth;
