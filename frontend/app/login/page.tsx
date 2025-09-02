'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import MatrixRain from '../components/matrixrain';

export default function LoginPage() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    router.push('/dashboard');
  };

  return (
    <div className="flex h-screen items-center justify-center font-mono">
      {/* Matrix rain background */}
      <MatrixRain />

      <form
        onSubmit={handleLogin}
        className="flex w-80 flex-col gap-4 rounded-xl border border-green-500 bg-black/80 p-6 shadow-[0_0_20px_#00ff00]"
      >
        <h1 className="mb-2 text-center text-3xl font-bold text-green-400 tracking-wider drop-shadow-[0_0_10px_#00ff00]">
          Gate of Things
        </h1>

        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="rounded-lg border border-green-500 bg-black p-2 text-green-400 placeholder-green-700 focus:outline-none focus:shadow-[0_0_10px_#00ff00]"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="rounded-lg border border-green-500 bg-black p-2 text-green-400 placeholder-green-700 focus:outline-none focus:shadow-[0_0_10px_#00ff00]"
        />

        <button
          type="submit"
          className="mt-2 rounded-lg border border-green-500 bg-black p-2 font-bold text-green-400 transition hover:bg-green-500 hover:text-black hover:shadow-[0_0_15px_#00ff00]"
        >
          Enter
        </button>
      </form>
    </div>
  );
}
