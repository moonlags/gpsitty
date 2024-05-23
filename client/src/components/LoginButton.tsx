import { createSignal } from "solid-js";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { AiOutlineUser } from "solid-icons/ai";
import axios from "axios";
import { showToast } from "./ui/toast";

function LoginButton() {
  const [email, setEmail] = createSignal<string>();
  const [password, setPassword] = createSignal<string>();

  function handleLogin() {
    axios
      .post(import.meta.env.VITE_BACKEND_HOST + "/auth/login", {
        email: email(),
        password: password(),
      })
      .then(() => {
        location.reload();
      })
      .catch(() => {
        showToast({
          title: "ERROR",
          description: "Wrong email or password",
          variant: "error",
        });
      });
  }

  function handleRegister() {
    axios
      .post(import.meta.env.VITE_BACKEND_HOST + "/auth/register", {
        email: email(),
        password: password(),
      })
      .then(() => {
        location.reload();
      })
      .catch(() => {
        showToast({
          title: "ERROR",
          description: "Email already registered",
          variant: "error",
        });
      });
  }

  return (
    <Popover>
      <PopoverTrigger>
        <Button class="flex flex-row gap-2 dark">
          <AiOutlineUser class="w-6 h-6" />
          Login
        </Button>
      </PopoverTrigger>
      <PopoverContent class="w-96 ml-2">
        <div class="grid gap-4">
          <div class="grid gap-2">
            <div class="grid grid-cols-3 items-center gap-4">
              <Label>Email</Label>
              <Input
                value={email()}
                onChange={(e) => setEmail(e.target.value)}
                class="col-span-2 h-8"
              />
            </div>
            <div class="grid grid-cols-3 items-center gap-4">
              <Label>Password</Label>
              <Input
                value={password()}
                onChange={(e) => setPassword(e.target.value)}
                class="col-span-2 h-8"
                type="password"
              />
            </div>
          </div>
          <div class="grid gap-2 mt-5">
            <Button class="dark border border-blue-100" onClick={handleLogin}>
              Login
            </Button>
            <p>---------------------------------------------------</p>
            <Button
              class="dark border border-blue-100"
              onClick={handleRegister}
            >
              Register
            </Button>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}

export default LoginButton;
