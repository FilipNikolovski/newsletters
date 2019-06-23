import React, { Component } from "react";
import { Router, Route, Switch } from "react-router-dom";
import { Box, Grommet } from "grommet";
import Landing from "./Landing";
import { AuthProvider } from "./Auth/AuthContext";
import Dashboard from "./Dashboard";
import Logout from "./Auth/Logout";
import VerifyEmail from "./VerifyEmail";
import ProtectedRoute from "./ProtectedRoute";
import history from "./history";

const theme = {
  global: {
    font: {
      family:
        "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol'",
      size: "14px",
      height: "20px"
    },
    colors: {
      background: "#F5F7F9",
      brand: "#6650AA"
    }
  },
  tabs: {
    header: {
      background: "white"
    }
  },
  tab: {
    color: "#888888",
    active: {
      color: "brand"
    },
    border: false
  },
  formField: {
    label: {
      color: "dark-3",
      size: "small",
      margin: { vertical: "0", top: "small", horizontal: "0" },
      weight: 600
    },
    border: false,
    borderColor: "#CACACA",
    margin: 0
  },
  button: {
    border: {
      radius: "5px",
      color: "#6650AA"
    },
    padding: {
      vertical: "7px",
      horizontal: "24px"
    },
    primary: {
      color: "#6650AA"
    },
    extend: props => {
      let extraStyles = "";
      if (props.primary) {
        extraStyles = `
            text-transform: uppercase;
          `;
      }
      return `
          color: white;
          font-size: 12px;
          font-weight: bold;
          border: 0px;
          border-radius:5px;
          ${extraStyles}
        `;
    }
  },
  anchor: {
    primary: {
      color: "#999999"
    },
    color: "#6650AA"
  }
};

class App extends Component {
  render() {
    return (
      <Grommet theme={theme} full>
        <Router history={history}>
          <AuthProvider>
            <Box flex background="background">
              <Switch>
                <ProtectedRoute path="/dashboard" component={Dashboard} />
                <Route path="/logout" component={Logout} />
                <Route path="/verify-email/:token" component={VerifyEmail} />
                <Route path="/" component={Landing} />
              </Switch>
            </Box>
          </AuthProvider>
        </Router>
      </Grommet>
    );
  }
}

export default App;
