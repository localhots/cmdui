import React, { Component } from 'react'
import './header.css'

export default class Header extends Component {
    render() {
        let user = this.props.user
        return (
            <header>
                <div className="header-container">
                    <div className="go-api">go-api<span>commands</span></div>
                    <div className="auth">
                        <div className="name">{user.name}</div>
                        <img src={user.picture} alt={"Picture of "+ user.name} />
                    </div>
                </div>
            </header>
        )
    }
}
