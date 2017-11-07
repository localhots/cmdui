import React, { Component } from 'react'
import { Link } from 'react-router-dom'

import Nav from '../blocks/nav.js'
import Output from '../blocks/output.js'
import Alert from '../blocks/alert.js'
import { api, httpPOST } from '../http.js'
import './command.css'

export default class Command extends Component {
    constructor(props) {
        super(props)
        this.state = {
            form: this.defaultForm(props.commands[props.cmd]),
            jobID: null,
            error: null,
        }
    }

    componentWillReceiveProps(props) {
        let form = this.defaultForm(props.commands[props.cmd])
        this.setState({form: form})
    }

    defaultForm(cmd) {
        if (!cmd || Object.keys(cmd).length === 0) {
            return {}
        }

        return {
            command: cmd.name,
            args: "",
            flags: cmd.flags.reduce((list, flag) => {
                list[flag.name] = flag.default
                return list
            }, {})
        }
    }

    render() {
        return (
            <div className="container">
                <Nav commands={this.props.commands} active={this.props.cmd}
                    query={this.props.query} onQueryChange={this.props.onQueryChange} />
                <main>
                    {this.renderForm()}
                    <Alert type="error" message={this.state.error} />
                    <Output jobID={this.state.jobID} />
                </main>
            </div>
        )
    }

    renderForm() {
        let cmd = this.props.commands[this.props.cmd]
        if (!cmd) {
            return null
        }

        return (
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
                    {cmd.flags.map(this.renderFlag.bind(this))}
                </ul>
                <input type="submit" value="Execute" id="submit-button"/>
                <Link to={"/cmd/"+ cmd.name +"/jobs"} className="history">History</Link>
            </form>
        )
    }

    renderFlag(flag) {
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
            var inputType = "string"
            let numericTypes = [
                "int", "int8", "int16", "int32", "int64",
                "uint", "uint8", "uint16", "uint32", "uint64"
            ]
            if (numericTypes.includes(flag.type)) {
                inputType = "number"
            }
            return (
                <li key={flag.name}>
                    <label htmlFor={"flag-"+flag.name}>{flag.name}</label>
                    <span className="descr">{flag.description}</span>
                    <input type={inputType} name={flagName(flag.name)} id={"flag-"+flag.name}
                        defaultValue={flag.default}
                        onChange={this.changeHandler.bind(this)} />
                </li>
            )
        }
    }

    submitHandler(event) {
        event.preventDefault()
        httpPOST(api("/exec"), formQuery(this.state.form),
            (response) => {
                let details = JSON.parse(response)
                this.setState({jobID: details.id, error: null})
            },
            (error) => {
                this.setState({jobID: null, error: error})
            },
        )
    }

    changeHandler(event) {
        let field = event.target
        let form = this.state.form
        if (field.id.indexOf("flag-") === 0) {
            var val = field.value
            if (field.type === "checkbox") {
                val = JSON.stringify(field.checked)
            }

            let name = field.id.substring(5)
            form.flags[name] = val

            this.setState({
                form: {
                    command: form.command,
                    args: form.args,
                    flags: form.flags
                }
            })
        } else if (field.name === "args") {
            this.setState({
                form: {
                    command: form.command,
                    args: field.value,
                    flags: form.flags
                }
            })
        }
    }
}

function formQuery(form) {
    var args = Object.keys(form.flags).map((name) => ["flags[" + name + "]", form.flags[name]])
    args.unshift(["command", form.command], ["args", form.args])
    let param = (pair) => pair[0] + "=" + encodeURIComponent(pair[1])
    return args.map(param).join("&")
}
