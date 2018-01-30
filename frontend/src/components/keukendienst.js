import React, { Component } from 'react'
import { Row, Preloader, Col, Input, Button, Icon, Table } from 'react-materialize'
import { getKeukendienst, setKeukendienst } from '../apis/keukendienst_api'

class KeukenDienst extends Component {
    content = {}

    constructor(props) {
        super(props)

        this.state = { keukenDienst:{}, loading: true, isSaving: false, }
        getKeukendienst().then(this.onData.bind(this))

        this.save = this.save.bind(this)
        this.addTask = this.addTask.bind(this)
        this.setContent = this.setContent.bind(this)
        this.deleteTask = this.deleteTask.bind(this)
    }

    onData(res) {
        this.setState({ loading: false, keukenDienst: res.data})
        this.content = res.data.content
        console.log(res.data)
    }

    save() {
        const keukenDienst = this.state.keukenDienst
        keukenDienst.content = this.content
        this.setState({ keukenDienst, isSaving: true })
        setKeukendienst(keukenDienst).then(this.onDoneSave.bind(this))
    }

    onDoneSave() {
        this.setState({ isSaving: false })
    }

    addTask() {
        const keukenDienst = this.state.keukenDienst
        if (!this.newTask.input.value) {
            return
        }
        keukenDienst.tasks.push(this.newTask.input.value)
        this.newTask.input.value = ""

        this.setState({ keukenDienst })
    }

    setContent(day,task,name) {
        if (!this.content[day]) {
            this.content[day] = {}
        }
        this.content[day][task] = name
    }

    deleteTask(task) {
        const keukenDienst = this.state.keukenDienst
        keukenDienst.tasks.splice(keukenDienst.tasks.indexOf(task),1)
        for (let key in this.content) {
            if (this.content.hasOwnProperty(key)) {
                delete this.content[key][task]
            }
        }
        this.setState({ keukenDienst })
    }

    render() {
        if (this.state.loading) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }

        const days = this.state.keukenDienst.days.map((i,j) => <th key={j}>{i}</th>)
        const rows = this.state.keukenDienst.tasks.map((task,j) => {
            return (
            <tr key={j}>
                <th>{task}</th>
                {this.state.keukenDienst.days.map((day, j) => <td key={j}><Input defaultValue={(this.state.keukenDienst.content[day] || {})[task]} onChange={(c) => this.setContent(day,task,c.target.value)} /></td>)}
                <td><a onClick={() => this.deleteTask(task)}><Icon right>delete</Icon></a></td>
            </tr>
            )
        })

        return <div>
            <Row><h2>Keuken Dienst</h2></Row>

            <Row>
                <Table>
	                <thead>
		                <tr>
                            <th></th>
                            {days}
                        </tr>
	                </thead>
                    <tbody>
                        {rows}
                        <tr>
                            <td>
                                <Input label="Nieuwe Taak" ref={(c) => this.newTask = c} validate/>
                                <Button waves='light' onClick={this.addTask} icon="add" floating/>
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