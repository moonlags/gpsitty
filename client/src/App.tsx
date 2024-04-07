import { useState } from "react";
import contentWarning from "/content_warning.webm";

function App() {
  const [count, setCount] = useState(0);

  const handleLogin = () => {
    window.location.href = "http://localhost:50731/auth/google";
  };

  return (
    <div className="h-screen w-full justify-center flex items-center flex-col gap-12">
      <video controls className="rounded-lg shadow-xl border-2 border-gray-500">
        <source src={contentWarning} type="video/webm" />
      </video>
      <div className="flex flex-row gap-10">
        <button
          className="bg-gray-300 px-10 py-3 rounded-md shadow-md font-bold hover:scale-110 duration-100"
          onClick={() => setCount((count) => count + 1)}
        >
          count is {count}
        </button>
        <button
          className="bg-gray-300 px-10 py-3 rounded-md shadow-md font-bold hover:scale-110 duration-100"
          onClick={handleLogin}
        >
          Login
        </button>
      </div>
    </div>
  );
}

export default App;
