import React, { Component } from 'react';
import './login.css';
import githubLogo from './github_logo.png';

export default class Login extends Component {
    componentDidMount() {
        let page = document.getElementsByClassName("login-page")[0];
        page.style.height = window.innerHeight + "px";
    }

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
