import React, { Component } from 'react'
import { Row, Preloader, Col, Input, Button, Icon, Table } from 'react-materialize'
import { getKeukendienst, setKeukendienst } from '../apis/keukendienst_api'
import DatePicker from 'react-datepicker'
import moment from "moment";
import 'react-datepicker/dist/react-datepicker.css';
import '../css/picker.css'

class KeukenDienst extends Component {

    constructor(props) {
        super(props)

        this.state = { keukenDienst:{}, loading: true, isSaving: false, }
        getKeukendienst().then(this.onData.bind(this))

        this.save = this.save.bind(this)
        this.addWeek = this.addWeek.bind(this)
        this.setFromTime = this.setFromTime.bind(this)
        this.setToTime = this.setToTime.bind(this)
        this.deleteTask = this.deleteTask.bind(this)
        this.setName = this.setName.bind(this)
    }

    onData(res) {
        for (let id in res.data) {
            res.data[id].from = moment(res.data[id].from)
            res.data[id].to = moment(res.data[id].to)
        }
        this.setState({ loading: false, keukenDienst: res.data})
        console.log(res.data)
    }

    save() {
        this.setState({ isSaving: true })
        setKeukendienst(this.state.keukenDienst).then(this.onDoneSave.bind(this))
    }

    onDoneSave() {
        this.setState({ isSaving: false })
    }

    setFromTime(i, content) {
        console.log(i, content)
        const keukenDienst = this.state.keukenDienst
        console.log(keukenDienst)
        keukenDienst[i].from = content
        console.log(keukenDienst)
        this.setState({ keukenDienst })
    }

    setToTime(i, content) {
        const keukenDienst = this.state.keukenDienst
        keukenDienst[i].to = content
        this.setState({ keukenDienst })
    }

    setName(i, j, name) {
        console.log(i, j, name)
        const keukenDienst = this.state.keukenDienst
        keukenDienst[i].names[j] = name
        this.setState({ keukenDienst })
    }

    addWeek() {
        const keukenDienst = this.state.keukenDienst
        keukenDienst.push({from: moment(new Date()), to: moment(new Date()), names:["", ""]})

        this.setState({ keukenDienst })
    }

    deleteTask(i) {
        const keukenDienst = this.state.keukenDienst
        keukenDienst.splice(i, 1)
        this.setState({ keukenDienst })
    }

    render() {
        if (this.state.loading) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }

        const rows = this.state.keukenDienst.map((task,j) => {
            return (
            <tr key={j}>
                <td><DatePicker label="Van" dateFormat="YYYY-MM-DD" selected={task.from} onChange={(c) => this.setFromTime(j, c)}/></td>
                <td><DatePicker label="Tot" dateFormat="YYYY-MM-DD" selected={task.to} onChange={(c) => this.setToTime(j, c)}/></td>
                <td><Input label="Naam 1" validate defaultValue={task.names[0]} onChange={(i, c) => this.setName(j, 0, c)}/></td>
                <td><Input label="Naam 2" validate defaultValue={task.names[1]} onChange={(i, c) => this.setName(j, 1, c)}/></td>
                <td><a onClick={() => this.deleteTask(j)}><Icon right>delete</Icon></a></td>
            </tr>
            )
        })

        return <div className="containerHeight">
            <Row><h2>Keuken Dienst</h2></Row>
            <Row>
                <Table>
	                <thead>
		                <tr>
                            <th>Van</th>
                            <th>Tot</th>
                            <th>Naam 1</th>
                            <th>Naam 2</th>
                            <th/>
                        </tr>
	                </thead>
                    <tbody>
                        {rows}
                        <tr>
                            <td>
                                <Button waves='light' onClick={this.addWeek} icon="add" floating/>
                            </td>
                        </tr>
                    </tbody>
                </Table>
            </Row>
            
            <Row>
                <Col s={2}><Button waves='light' disabled={this.state.isSaving} onClick={this.save}>Opslaan<Icon left>save</Icon></Button></Col>
            </Row>
        </div>
    }
}

export default KeukenDienst