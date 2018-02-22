import React, { Component } from 'react'
import { Row, Preloader, Col, Input, Button, Icon } from 'react-materialize'
import { withRouter } from "react-router-dom";

import { addTickerItem } from '../apis/ticker_api'

class TickerEdit extends Component {
    constructor(props) {
        super(props)


        this.state = { isSaving: false }
        this.item = { setup: "", metric: "", interval: "10s", back: "10s" }

        this.save = this.save.bind(this)
        this.doneSaving = this.doneSaving.bind(this)
    }


    save() {
        this.setState({ isSaving: false })
        addTickerItem({ setup: this.item.setup, metric: this.item.metric, interval: this.item.interval, back: this.item.back }).then(this.doneSaving)
    }

    doneSaving() {
        this.setState({ isSaving: false })
        this.props.history.push("/ticker");
    }

    render() {
        return <div className="marginBottom containerHeightNews">
            <Row><h2>Ticker item toevoegen</h2></Row>
            <Row>
	            <Input s={12} label="Opstelling" validate defaultValue={this.item.setup} onChange={(c, value) => this.item.setup = value} />
            </Row>
            <Row>
	            <Input s={12} label="Meting" validate defaultValue={this.item.metric} onChange={(c, value) => this.item.metric = value} />
            </Row>
            <Row>
	            <Input s={4} label="Interval" validate defaultValue={this.item.interval} onChange={(c, value) => this.item.interval = value} />
            </Row>
            <Row>
	            <Input s={4} label="Terugname" validate defaultValue={this.item.back} onChange={(c, value) => this.item.back = value} />
            </Row>
            
            
            <Row>
                <Col s={2}><Button waves='light' disabled={this.state.isSaving} onClick={this.save}>Opslaan<Icon left>save</Icon></Button></Col>
            </Row>
        </div>
    }
}

export default withRouter(TickerEdit)