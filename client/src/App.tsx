import { useEffect, useState } from "react";
import axios from "axios";
import { IUser } from "./components/UserInfo";
import NavBar from "./components/Navbar";
import DevicesView from "./components/DevicesView";
import IntorductionPage from "./components/IntroductionPage";

function App() {
  const [user, setUser] = useState<IUser>();

  useEffect(() => {
    axios
      .get("http://" + import.meta.env.VITE_BACKEND_HOST + "/api/v1/session", {
        withCredentials: true,
      })
      .then((response) => {
        setUser(response.data);
      })
      .catch((err) => {
        console.error(err);
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
