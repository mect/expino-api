import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route } from 'react-router-dom';

import Header from './components/header'
import Footer from './components/footer'
import Home from './components/home'
import NewsOverview from './components/newsoverview'
import NewsEdit from './components/newsedit'
import FeatureSlides from './components/featureslides'
import KeukenDienst from './components/keukendienst'
import TickerItems from './components/tickeritems';
import TickerEdit from './components/tickeredit'
import GraphItems from './components/graphitems'
import GraphEdit from './components/graphedit'

const App = function () {
    return (
        <div>
        <Router>
            <div>
                <Header />
                <div className="container">
                    <Route path="/" exact component={Home} />
                    <Route path="/news" exact component={NewsOverview} />
                    <Route path="/news/edit/:id?" component={NewsEdit} />
                    <Route path="/featureslides" exact component={FeatureSlides}/>
                    <Route path="/keukendienst" exact component={KeukenDienst}/>
                    <Route path="/ticker" exact component={TickerItems} />
                    <Route path="/ticker/new" exact component={TickerEdit} />
                    <Route path="/graphs" exact component={GraphItems} />
                    <Route path="/graphs/new" exact component={GraphEdit} />
                </div>
                <Footer />
            </div>
        </Router>
    </div>
    );
};

ReactDOM.render(
    <App />,
    document.querySelector("#container")
);