import React, { Component } from 'react';

import './login.css';
import githubLogo from './github_logo.png';

export default class Login extends Component {
    render() {
        return (
            <div className="login-page">
                <a className="button" href="/api/auth/login">
                    <img src={githubLogo} alt="Login with GitHub" />
                    <div>Login with GitHub</div>
                </a>
            </div>
        );
    }
}
