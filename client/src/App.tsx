import { useState } from "react";
import "./App.css";

function App() {
  const [count, setCount] = useState(0);

  const handleLogin = () => {
    window.location.href = "http://localhost:50731/auth/google";
  };

  return (
    <>
      <div>
        <iframe
          className="m-10 w-[500px] h-[400px]"
          src="https://open.spotify.com/embed/playlist/3PUMIm7UdUzJimwmGuXGcV?utm_source=generator&theme=0"
          allowFullScreen={false}
          allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
          loading="lazy"
        ></iframe>
      </div>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
      </div>
      <code>dasdsad</code>
      <div>
        <button onClick={handleLogin}>Sing in with google</button>
      </div>
    </>
  );
}

export default App;
