import React, { Component } from 'react';
import { Row, Collection, CollectionItem, Button , Icon} from 'react-materialize';
import { Link } from 'react-router-dom'
import { getTickerItems, deleteTickerItem } from '../apis/ticker_api' 
import _ from 'lodash/collection'

class TickerItems extends Component {
    constructor(props) {
        super(props)
        this.state = { items: [] }

        getTickerItems().then(this.gotItems.bind(this))
        this.deleteItem = this.deleteItem.bind(this)
    }

    gotItems(result) {
        this.setState({ items: result.data })
    }

    deleteItem(id) {
        this.setState({ items: _.reject(this.state.items, (i) => i.id === id) })
        deleteTickerItem(id)
    }

    render() {
        const items = this.state.items.map(i => <CollectionItem>
                {i.setup}: {i.metric}
                <a onClick={() => this.deleteItem(i.id)}><Icon right>delete</Icon></a>
            </CollectionItem>)
        return <div className="containerHeight">
                <Row>
                <Collection header='Ticker items'>
                    {items}
                </Collection>
                <Link to="/ticker/new"><Button floating fab='vertical' icon='add' className='red' large style={{bottom: '45px', right: '24px'}}/></Link>
                </Row>
                </div>
    }
}

export default TickerItems