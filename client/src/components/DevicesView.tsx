import axios from "axios";
import { createEffect, createSignal } from "solid-js";
import { showToast } from "./ui/toast";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import DeviceCard, { IDevice } from "./DeviceCard";

function DevicesView() {
  const [devices, setDevices] = createSignal<[IDevice]>();
  const [imei, setImei] = createSignal<string>();

  createEffect(() => {
    axios
      .get(import.meta.env.VITE_BACKEND_HOST + "/api/v1/devices")
      .then((response) => {
        console.log(response);
        setDevices(response.data);
      })
      .catch(() => {
        showToast({
          title: "ERROR",
          description: "Something went wrong while getting your devices",
          variant: "error",
        });
      });
  });

  function handleLink() {
    axios
      .get(import.meta.env.VITE_BACKEND_HOST + "/api/v1/link/" + imei())
      .then(() => {
        location.reload();
      })
      .catch(() => {
        showToast({
          title: "ERROR",
          description: "Something went wrong, check provided IMEI",
          variant: "error",
        });
      });
  }

  return (
    <div class="flex flex-col gap-20 w-full">
      <p class="text-3xl font-semibold">Your Devices:</p>
      <div class="px-20 grid sm:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6">
        {devices()?.map((device) => (
          <DeviceCard device={device} />
        ))}
      </div>
      <div class="flex flex-col items-start gap-2">
        <Popover>
          <PopoverTrigger class="text-blue-400 text-lg underline cursor-pointer">
            Can not find your device?
          </PopoverTrigger>
          <PopoverContent class="w-96 ml-2">
            <div class="grid gap-4">
              <div class="space-y-2">
                <h4 class="font-medium leading-none">Link It!</h4>
              </div>
              <div class="grid gap-2">
                <div class="grid grid-cols-3 items-center gap-4">
                  <Label>IMEI</Label>
                  <Input
                    value={imei()}
                    onChange={(e) => setImei(e.target.value)}
                    class="col-span-2 h-8"
                  />
                </div>
              </div>
              <Button class="dark border border-blue-100" onClick={handleLink}>
                Link
              </Button>
            </div>
          </PopoverContent>
        </Popover>
      </div>
    </div>
  );
}

export default DevicesView;
