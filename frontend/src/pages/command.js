import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import Nav from '../blocks/nav.js';
import Output from '../blocks/output.js';
import { api, httpPOST } from '../http.js';
import './command.css';

export default class Command extends Component {
    constructor(props) {
        super(props);

        this.state = {
            form: this.defaultForm(props.commands[props.cmd]),
            jobID: null,
            error: null,
        };
    }

    componentWillReceiveProps(props) {
        this.setState({form: this.defaultForm(props.commands[props.cmd])});
    }

    defaultForm(cmd) {
        if (!cmd || Object.keys(cmd).length === 0) {
            return {};
        }

        var form = {
            command: cmd.name,
            args: "",
            flags: {}
        }
        cmd.flags.forEach((flag) => {
            form.flags[flag.name] = flag.default;
        });
        return form;
    }

    render() {
        let cmd = this.props.commands[this.props.cmd];
        if (!cmd) {
            return null;
        }

        return (
            <div className="container">
                <Nav commands={this.props.commands} active={this.props.cmd}
                    query={this.props.query} onQueryChange={this.props.onQueryChange} />
                <main>
                    <form method="post" action="/api/exec" onSubmit={this.submitHandler.bind(this)} ref="form">
                        <div className="command">
                            <div className="name">{cmd.name}</div>
                            <div className="descr">{cmd.description}</div>
                            <input type="hidden" name="command" defaultValue={cmd.name} />
                        </div>
                        <ul className="fields">
                            <li>
                                <div className="descr">Command arguments</div>
                                <input type="text" name="args" placeholder={cmd.args_placeholder}
                                    onChange={this.changeHandler.bind(this)} />
                            </li>
                            {cmd.flags.map((flag) => { return this.input(flag); })}
                        </ul>
                        <input type="submit" value="Execute" id="submit-button"/>
                        <Link to={"/cmd/"+ cmd.name +"/jobs"} className="history">History</Link>
                    </form>
                    {this.renderError()}
                    <Output cmd={cmd} jobID={this.state.jobID} key={this.state.jobID} />
                </main>
            </div>
        );
    }

    renderError() {
        if (this.state.error === null) {
            return null;
        }
        return (
            <div className="error">{this.state.error}</div>
        );
    }

    submitHandler(event) {
        event.preventDefault();

        let f = this.state.form;
        var args = [];
        args.push(["command", f.command]);
        args.push(["args", f.args]);
        for (let name of Object.keys(f.flags)) {
            args.push(["flags[" + name + "]", f.flags[name]]);
        }
        let formQuery = args.map((pair) => {
            return pair[0] + "=" + encodeURIComponent(pair[1]);
        }).join("&");

        httpPOST(api("/exec"), formQuery,
            (response) => {
                let details = JSON.parse(response);
                this.setState({jobID: details.id, error: null});
            },
            (error) => {
                this.setState({jobID: null, error: error});
            },
        );
    }

    changeHandler(event) {
        if (event.target.id.indexOf("flag-") === 0) {
            var val = event.target.value;
            if (event.target.type === "checkbox") {
                val = JSON.stringify(event.target.checked);
            }

            var flags = this.state.form.flags;
            let name = event.target.id.substring(5);
            flags[name] = val;

            this.setState({
                form: {
                    command: this.state.form.command,
                    args: this.state.form.args,
                    flags: flags
                }
            });
        } else if (event.target.name === "args") {
            this.setState({
                form: {
                    command: this.state.form.command,
                    args: event.target.value,
                    flags: this.state.form.flags
                }
            });
        }
    }

    input(flag) {
        let flagName = (name) => "flags[" + name + "]"
        if (flag.type === "bool") {
            return (
                <li key={flag.name}>
                    <input type="checkbox" name={flagName(flag.name)} id={"flag-"+flag.name}
                        defaultChecked={(flag.default === "true")} value="true"
                        onClick={this.changeHandler.bind(this)} />
                    <label htmlFor={"flag-"+flag.name}>{flag.name}</label>
                    <span className="descr">{flag.description}</span>
                </li>
            )
        } else {
            var inputType = "string";
            let numericTypes = [
                "int", "int8", "int16", "int32", "int64",
                "uint", "uint8", "uint16", "uint32", "uint64"
            ];
            if (numericTypes.includes(flag.type)) {
                inputType = "number";
            }
            return (
                <li key={flag.name}>
                    <label htmlFor={"flags-"+flag.name}>{flag.name}</label>
                    <span className="descr">{flag.description}</span>
                    <input type={inputType} name={flagName(flag.name)} id={"flag-"+flag.name}
                        defaultValue={flag.default}
                        onChange={this.changeHandler.bind(this)} />
                </li>
            )
        }
    }
}
