import {useEffect, useState} from "react";

export default function App() {
    const [status, setStatus] = useState();
    const [installState, setInstallState] = useState(0);

    window.runtime.EventsOn("status", setStatus)
    window.runtime.EventsOn("setInstallState", setInstallState)

    useEffect(() => {
        if (1 === installState) {
            window.runtime.EventsEmit("install");
        }
    });

    return (
        <>
            <Instructions installState={installState}/>
            <InstallBtn installState={installState}  setInstallState={setInstallState}/>
            <p>{status}</p>
        </>
    )
}


function Instructions({installState}) {
    if (0 < installState) {
        return
    }

    return (
        <div>
            <p>Click "install" to load the Lume Web extension into Chrome and Brave</p>
            <p style={{color:"red"}}>Please <b>ensure</b> your Chrome and/or Brave browser is completely closed before proceeding!</p>
        </div>
    );
}

function InstallBtn({installState,setInstallState}) {
    if (0 < installState) {
        return
    }

    return <button onClick={() => setInstallState(1)}>Install</button>;
}