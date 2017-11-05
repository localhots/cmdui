import React, { Component } from 'react'

export default class Timestamp extends Component {
    constructor(props) {
        super(props)
        this.state = {timer: null, text: null}
    }

    componentDidMount() {
        if (this.props.relative) {
            this.setState({
                text: this.relativeDate(),
                timer: setInterval(this.setRelativeDate.bind(this), 1000)
            })
        } else {
            this.setState({
                text: null,
                timer: null
            })
        }
    }

    componentWillUnmount() {
        if (this.state.timer !== null) {
            clearInterval(this.state.timer)
            this.setState({timer: null, text: null})
        }
    }

    relativeDate(date = this.props.date) {
        return this.timeSince(new Date(date)) + " ago"
    }

    setRelativeDate(date = this.props.date) {
        this.setState({text: this.relativeDate()})
    }

    render() {
        return (
            <span title={this.formatTitle()}>{this.formatText()}</span>
        )
    }

    formatText() {
        if (this.props.date) {
            if (this.props.relative) {
                return this.state.text
            } else {
                return this.formatDate(new Date(this.props.date))
            }
        } else if (this.props.from && this.props.until) {
            return this.timeSince(new Date(this.props.from), new Date(this.props.until))
        } else {
            return "—"
        }
    }

    formatTitle() {
        if (this.props.relative) {
            return this.formatDate(new Date(this.props.date))
        } else if (this.props.until) {
            return this.formatDate(new Date(this.props.until))
        } else {
            return ""
        }
    }

    formatDate(date) {
        let leftPad = (n) => n > 9 ? n : "0"+ n // I wish I was a package
        let monthNames = [
            "Jan", "Feb", "Mar", "Apr", "May", "Jun",
            "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
        ]
        let y = date.getFullYear(),
            m = date.getMonth(),
            d = date.getDate(),
            h = leftPad(date.getHours()),
            i = leftPad(date.getMinutes()),
            s = leftPad(date.getSeconds())
        return [d, monthNames[m], y, "at", [h, i, s].join(":")].join(" ")
    }

    timeSince(timeStamp, until = new Date()) {
        let secondsPast = (until.getTime() - timeStamp.getTime())/1000
        if (secondsPast < 60) {
            return parseInt(secondsPast, 10) + "s"
        } else if (secondsPast < 3600) {
            return parseInt(secondsPast/60, 10) + "m"
        } else if (secondsPast <= 86400) {
            return parseInt(secondsPast/3600, 10) + "h"
        } else {
            let day = timeStamp.getDate()
            let month = timeStamp.toDateString().match(/ [a-zA-Z]*/)[0].replace(" ", "")
            let year = timeStamp.getFullYear() === until.getFullYear()
                ? ""
                : " "+ timeStamp.getFullYear()
            return day +" "+ month + year
        }
    }
}
