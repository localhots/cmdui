import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom'

import Loading from './pages/loading.js';
import Login from './pages/login.js';
import SelectCommand from './pages/select_command.js';
import Command from './pages/command.js';
import History from './pages/history.js';
import Job from './pages/job.js';
import Header from './blocks/header.js';
import { api, httpGET } from './http.js';
import './app.css';

export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authorized: null,
            user: null,
            commands: null,
            commands_query: ""
        };
        this.loadAuth();
        this.loadCommands();
    }

    loadAuth() {
        httpGET(api("/auth/session"),
            (status, body) => {
                if (status === 200) {
                    this.setState({
                        authorized: true,
                        user: JSON.parse(body)
                    });
                } else {
                    this.setState({authorized: false});
                }
            },
            (error) => {
                console.log("Failed to load auth details:", error);
                this.setState({authorized: false});
            }
        );
    }

    loadCommands() {
        httpGET(api("/commands"),
            (status, body) => {
                if (status === 200) {
                    let list = JSON.parse(body);
                    var hash = {};
                    for (var cmd of list) {
                        hash[cmd.name] = cmd;
                    }
                    this.setState({commands: hash});
                }
            },
            (error) => {
                this.setState({commands: {}});
                console.log("Failed to load commands:", error);
            }
        );
    }

    render() {
        if (this.state.authorized === null || this.state.commands === null) {
            return <Loading />;
        }
        if (this.state.authorized === false) {
            return <Login />;
        }

        return (
            <Router>
                <div className="router-container">
                    <Header user={this.state.user} />
                    <Route exact path="/" render={props => (
                        <SelectCommand commands={this.state.commands}
                            onQueryChange={this.queryHandler.bind(this)} query={this.state.commands_query} />
                    )}/>
                    {/* Command */}
                    <Route exact path={"/cmd/:name"} render={props => (
                        <Command
                            cmd={props.match.params.name} commands={this.state.commands}
                            query={this.state.commands_query} onQueryChange={this.queryHandler.bind(this)} />
                    )}/>
                    {/* Logs */}
                    <Route exact path="/cmd/:name/jobs/:jobID" render={props => (
                        <Job jobID={props.match.params.jobID}
                            cmd={props.match.params.name} commands={this.state.commands}
                            query={this.state.commands_query} onQueryChange={this.queryHandler.bind(this)} />
                    )}/>
                    {/* History */}
                    <Route exact path="/cmd/:name/jobs" render={props => (
                        <History cmd={props.match.params.name} commands={this.state.commands}
                            query={this.state.commands_query} onQueryChange={this.queryHandler.bind(this)} />
                    )}/>
                    <Route exact path="/users/:id/jobs" render={props => (
                        <History userID={props.match.params.id}
                            commands={this.state.commands}
                            query={this.state.commands_query} onQueryChange={this.queryHandler.bind(this)} />
                    )}/>
                    <Route exact path="/jobs" render={props => (
                        <History
                            commands={this.state.commands}
                            query={this.state.commands_query} onQueryChange={this.queryHandler.bind(this)} />
                    )}/>
                </div>
            </Router>
        );
    }

    queryHandler(event) {
        this.setState({
            commands_query: event.target.value,
        });
    }
}
