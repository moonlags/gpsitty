import { useEffect, useState } from "react";
import axios from "axios";
import { User } from "./components/UserInfo";
import NavBar from "./components/Navbar";
import DevicesView from "./components/DevicesView";
import IntorductionPage from "./components/IntroductionPage";

function App() {
  const [user, setUser] = useState<User>();

  useEffect(() => {
    axios
      .get("http://localhost:50731/v1/session", { withCredentials: true })
      .then((response) => {
        if (response.status === 200) {
          setUser(response.data);
        }
        // todo: else toast
      })
      .catch((err) => {
        console.error(err);
        // todo: toast
      });
  }, []);

  return (
    <main className="flex h-screen flex-col">
      <NavBar {...user} />
      <div className="flex p-10">
        {user ? <DevicesView /> : <IntorductionPage />}
      </div>
    </main>
  );
}

export default App;
