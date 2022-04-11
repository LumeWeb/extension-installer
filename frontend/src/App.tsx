import {useEffect, useState} from "react";
// @ts-ignore
import IconSetting from "./assets/images/icon-setting.png";
import "./assets/css/style.css";

export default function App() {
    const [status, setStatus] = useState();
    const [installState, setInstallState] = useState(0);

    window.runtime.EventsOn("status", setStatus)
    window.runtime.EventsOn("setInstallState", setInstallState)
    return (
        <div className="content-container">
            <Icon installState={installState} />
            <Header installState={installState}/>
            <Instructions installState={installState}/>
            <InstallBtn installState={installState} setInstallState={setInstallState}/>
            {installState > 0 && <p>{status}</p>}
        </div>
    )
}

function Icon({installState}) {
    let classes = ['header-icon'];
    if (1 === installState) {
        classes.push("rotating");
    }

    return <img className={classes.join(" ")} src={IconSetting}/>
}

function Header({installState}) {
    let message = "Lume Web Extension Installer";

    if (1 == installState) {
        message = "Installing Extension";
    }

    if (2 == installState) {
        message = "Extension installed successfully";
    }

    return <h1>{message}</h1>;
}


function Instructions({installState}) {
    if (0 < installState) {
        return
    }

    return (
        <>
            <div className="paragraph">
                <p>Click install to load the Lume Web Extension into Chrome and Brave.</p>
            </div>
            <div className="error-message">

                <div className="error-message-content">
                    <span>Please ensure your Chrome and/or Brave browser is <strong>completely closed</strong> before proceeding.</span>
                </div>
            </div>
        </>
    );
}

function InstallBtn({installState, setInstallState}) {
    if (0 < installState) {
        return
    }

    const startInstall = () => {
        window.runtime.EventsEmit("install");
    };

    return (
        <div className="action-buttons">
            <button className="btn" onClick={() => startInstall()}>Install extension</button>
        </div>
    );
}