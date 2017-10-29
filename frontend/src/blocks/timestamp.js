import React, { Component } from 'react';

export default class Timestamp extends Component {
    constructor(props) {
        super(props);
        this.state = {timer: null, text: null};
    }

    componentDidMount() {
        if (this.props.relative) {
            this.setState({
                text: this.relativeDate(),
                timer: setInterval(this.setRelativeDate.bind(this), 1000)
            });
        } else {
            this.setState({
                text: null,
                timer: null
            });
        }
    }

    componentWillUnmount() {
        if (this.state.timer !== null) {
            clearInterval(this.state.timer);
            this.setState({timer: null, text: null});
        }
    }

    relativeDate(date = this.props.date) {
        return this.timeSince(new Date(date)) + " ago";
    }

    setRelativeDate(date = this.props.date) {
        this.setState({text: this.relativeDate()});
    }

    render() {
        var text;
        if (this.props.date !== undefined) {
            if (this.props.relative) {
                text = this.state.text;
            } else {
                text = this.formatDate(new Date(this.props.date));
            }
        } else if (this.props.from !== undefined && this.props.until !== undefined) {
            text = this.timeSince(new Date(this.props.from), new Date(this.props.until));
        } else {
            text = "&mdash;";
        }

        var title = "";
        if (this.props.relative) {
            title = this.formatDate(new Date(this.props.date));
        } else if (this.props.until !== undefined) {
            title = this.formatDate(new Date(this.props.until));
        }

        return (
            <span title={title}>{text}</span>
        );
    }

    formatDate(date) {
        let leftPad = (n) => n > 9 ? n : "0"+ n; // I wish I was a package
        let monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        let y = date.getFullYear(),
            m = date.getMonth(),
            d = date.getDate(),
            h = leftPad(date.getHours()),
            i = leftPad(date.getMinutes()),
            s = leftPad(date.getSeconds());
        return [d, monthNames[m], y, "at", [h, i, s].join(":")].join(" ");
    }

    timeSince(timeStamp, until = new Date()) {
        let secondsPast = (until.getTime() - timeStamp.getTime())/1000;
        if (secondsPast < 60) {
            return parseInt(secondsPast, 10) + 's';
        } else if (secondsPast < 3600) {
            return parseInt(secondsPast/60, 10) + 'm';
        } else if (secondsPast <= 86400) {
            return parseInt(secondsPast/3600, 10) + 'h';
        } else {
            let day = timeStamp.getDate();
            let month = timeStamp.toDateString().match(/ [a-zA-Z]*/)[0].replace(" ","");
            let year = timeStamp.getFullYear() === until.getFullYear()
                ? ""
                :  " "+timeStamp.getFullYear();
            return day + " " + month + year;
        }
    }
}
