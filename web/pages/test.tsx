import React from 'react';
import { connect } from 'react-redux';
import { IUserForm, registerAction } from '../states/auth/actions';
import { IAuth } from '../states/auth/state';
import { IStore } from 'states/type';

interface IProps {
  register: (form: IUserForm) => void;
  auth: IAuth;
}

const test: React.FC<IProps> = ({ auth, register }) => {
  return (
    <div>
      {auth.authenticated && auth.user.email + ' is registered 🌚!'}
      <br></br>
      <button
        onClick={() =>
          register({ email: 'Itworks@gmail.com', password: 'myepicpassword' })
        }
      >
        REGISTER!!!!!!!!🌝
      </button>
    </div>
  );
};

const mapStateToProps = (state: IStore) => ({
  auth: state.auth
});

export default connect(mapStateToProps, { register: registerAction })(test);
