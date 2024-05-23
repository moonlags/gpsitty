import { createEffect, createSignal } from "solid-js";
import { IUser } from "./components/UserInfo";
import axios from "axios";
import NavBar from "./components/Navbar";
import IntroductionPage from "./components/IntroductionPage";
import DevicesView from "./components/DevicesView";

function App() {
  const [user, setUser] = createSignal<IUser>();

  createEffect(() => {
    axios
      .get(import.meta.env.VITE_BACKEND_HOST + "/api/v1/session")
      .then((response) => {
        console.log(response);
        setUser(response.data);
      })
      .catch((err) => {
        console.error(err);
      });
  });

  return (
    <main class="flex h-screen flex-col">
      <NavBar user={user()} />
      <div class="flex p-10">
        {user()?.ID ? <DevicesView /> : <IntroductionPage />}
      </div>
    </main>
  );
}

export default App;
