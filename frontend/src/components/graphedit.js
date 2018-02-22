import React, { Component } from 'react'
import { Row, Col, Input, Button, Icon } from 'react-materialize'
import { withRouter } from "react-router-dom";

import { addGraphItem } from '../apis/graphs_api'

class GraphEdit extends Component {
    constructor(props) {
        super(props)


        this.state = { isSaving: false }
        this.item = { name: "" , url: "" }

        this.save = this.save.bind(this)
        this.doneSaving = this.doneSaving.bind(this)
    }


    save() {
        this.setState({ isSaving: false })
        addGraphItem(this.item).then(this.doneSaving)
    }

    doneSaving() {
        this.setState({ isSaving: false })
        this.props.history.push("/graphs");
    }

    render() {
        return <div className="marginBottom containerHeightNews">
            <Row><h2>Grafana item toevoegen</h2></Row>
            <Row>
	            <Input s={12} label="Naam" validate defaultValue={this.item.name} onChange={(c, value) => this.item.name = value} />
            </Row>
            <Row>
	            <Input s={12} label="URL" validate defaultValue={this.item.url} onChange={(c, value) => this.item.url = value} />
            </Row>

            <Row>
                <Col s={2}><Button waves='light' disabled={this.state.isSaving} onClick={this.save}>Opslaan<Icon left>save</Icon></Button></Col>
            </Row>
        </div>
    }
}

export default withRouter(GraphEdit)