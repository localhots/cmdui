import React, { Component } from 'react';

import { api, httpGET } from '../http.js';
import './user.css';

export default class User extends Component {
    constructor(props) {
        super(props);
        this.state = {user: undefined};
    }
    componentDidMount() {
        if (this.props.id === undefined || this.props.id === null) {
            return;
        }
        httpGET(api("/users/"+ this.props.id),
            (status, body) => {
                this.setState({user: JSON.parse(body)});
            },
            (error) => {
                console.log("Failed to load user details:", error);
            }
        );
    }

    render() {
        if (this.props.id === undefined || this.props.id === null || this.state.user === undefined) {
            return null;
        } else if (this.state.user === null) {
            let shortID = this.props.id.substring(0, 8);
            return (
                <div className="user-details">
                    User
                    <div className="user-id">{shortID}</div>
                </div>
            );
        } else {
            let details = this.state.user;
            return (
                <div className="user-details">
                    <img src={details.picture} alt={details.name +" picture"} />
                    <div className="name">{details.name}</div>
                </div>
            );
        }
    }
}
