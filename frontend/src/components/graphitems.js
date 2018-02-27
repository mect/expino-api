import React, { Component } from 'react';
import { Row, Collection, CollectionItem, Button , Icon} from 'react-materialize';
import { Link } from 'react-router-dom'
import { getGraphItems, deleteGraphItem } from '../apis/graphs_api' 
import _ from 'lodash/collection'

class GraphItems extends Component {
    constructor(props) {
        super(props)
        this.state = { items: [] }

        getGraphItems().then(this.gotItems.bind(this))
        this.deleteItem = this.deleteItem.bind(this)
    }

    gotItems(result) {
        this.setState({ items: result.data })
    }

    deleteItem(id) {
        this.setState({ items: _.reject(this.state.items, (i) => i.id === id) })
        deleteGraphItem(id)
    }

    render() {
        const items = this.state.items.map(i => <CollectionItem>
                {i.name}
                <a onClick={() => this.deleteItem(i.id)}><Icon right>delete</Icon></a>
            </CollectionItem>)
        return <div className="containerHeight">
                <Row>
                    <Button waves='light' node='a' href={`http://${window.location.hostname}:8000`} target="_blank"> Ga naar Grafana </Button>
                </Row>
                <Row>
                <Collection header='Grafana items'>
                    {items}
                </Collection>
                <Link to="/graphs/new"><Button floating fab='vertical' icon='add' className='red' large style={{bottom: '45px', right: '24px'}}/></Link>
                </Row>
                </div>
    }
}

export default GraphItems