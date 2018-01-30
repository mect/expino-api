import React, { Component } from 'react'
import { Row, Preloader, Col, Input, Button, Icon } from 'react-materialize'
import { getEnabledFeatureSlides, setEnabledFeatureSlides } from '../apis/settings_api'

const components = ["traffic", "social", "trains", "keukendienst"]
const componentNames = {
    "traffic": "Verkeer",
    "social": "Sociale Media",
    "trains": "Treinuren",
    "keukendienst": "Keuken dienst",
}

class FeatureSlides extends Component {
    checkboxes = {}

    constructor(props) {
        super(props)

        this.state = { componentsEnabled:[], loading: true, isSaving: false, }
        getEnabledFeatureSlides().then(this.onFeatureSlides.bind(this))

        this.save = this.save.bind(this)
    }

    onFeatureSlides(res) {
        this.setState({ loading: false, componentsEnabled: res.data})
        console.log(res.data)
    }

    save() {
        const enabled = [] 
        for (let key in this.checkboxes) {
            console.log(this.checkboxes[key])
            if (this.checkboxes.hasOwnProperty(key) && this.checkboxes[key].state.checked) {
                enabled.push(key)
            }
        }
        console.log(enabled)
        this.setState({ isSaving: true })
        setEnabledFeatureSlides(enabled).then(this.onDoneSave.bind(this))
    }

    onDoneSave() {
        this.setState({ isSaving: false })
    }

    render() {
        if (this.state.loading) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }

        const options = components.map((i,k) => <Row key={k}><Input name={i} type='checkbox' value={i} label={componentNames[i]}  ref={(c) => this.checkboxes[i] = c} checked={this.state.componentsEnabled.indexOf(i) > -1 ? true : false} /></Row>)

        return <div>
            <Row><h2>Slides</h2></Row>
            {options}
            <Row>
                <Col s={2}><Button waves='light' disabled={this.state.isSaving} onClick={this.save}>Opslaan<Icon left>save</Icon></Button></Col>
            </Row>
        </div>
    }
}

export default FeatureSlides