import { ClerkLoaded } from "@clerk/nextjs";
import { HeroUIProvider } from "@heroui/react";

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ClerkLoaded>
      <HeroUIProvider>{children}</HeroUIProvider>
    </ClerkLoaded>
  );
}
