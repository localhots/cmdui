import React, { Component } from 'react';

import './alert.css';

export default class Alert extends Component {
    render() {
        if (!this.props.message) {
            return null
        }
        return (
            <div className={"alert "+ this.props.type}>{this.props.message}</div>
        )
    }
}
