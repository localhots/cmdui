import React, { Component } from 'react';

import './header.css';

class Header extends Component {
    render() {
        return (
            <header>
                <div className="header-container">
                    <div className="go-api">go-api<span>commands</span></div>
                    <div className="auth">
                        <div className="name">{this.props.user.name}</div>
                        <img src={this.props.user.picture} alt="" />
                    </div>
                </div>
            </header>
        );
    }
}

export default Header;
