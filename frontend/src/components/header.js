import React from 'react';
import { Navbar, Row } from 'react-materialize';
import { NavLink } from 'react-router-dom'

export default () => {
    return <Row>
        <Navbar brand='Expino Kiosk' right>
            <li><NavLink to="/">Home</NavLink></li>
            <li><NavLink to="/news">Nieuws</NavLink></li>
            <li><NavLink to="/featureslides">Slides</NavLink></li>
        </Navbar>
    </Row>
}