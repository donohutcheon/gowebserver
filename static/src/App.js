import React from 'react';
import Dashboard from "./components/Dashboard/Dashboard";
import {
    Switch,
    BrowserRouter as Router
} from "react-router-dom";
import Login from "./components/Login";
import AuthContextProvider from './contexts/AuthContext';
import PublicRoute from "./components/routes/PublicRoute";
import PrivateRoute from "./components/routes/PrivateRoute";
import SignUp from "./components/SignUp";
import About from "./pages/About";

function App() {
    return (
        <div className="App">
            <AuthContextProvider>
                <Router>
                    <Switch>
                        <PublicRoute restricted={false} path="/login" component={Login} exact/>
                        <PublicRoute restricted={false} path="/sign-up" component={SignUp} exact/>
                        <PublicRoute restricted={false} path="/about" component={About} exact/>
                        <PrivateRoute restricted={true} path="/dashboard" component={Dashboard} exact/>
                        <PrivateRoute restricted={true} path="/" component={Dashboard} exact/>
                    </Switch>
                </Router>
            </AuthContextProvider>
        </div>
    );
}

export default App;
