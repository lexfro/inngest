import Link from 'next/link';

type IntegrationsLayoutProps = {
  children: React.ReactNode;
};

export default function IntegrationsLayout({ children }: IntegrationsLayoutProps) {
  return (
    <div className="min-height-[300px] flex divide-x divide-slate-100">
      <nav className="w-60 shrink-0 p-8">
        <ul>
          <li>
            <Link
              className="block w-full rounded-md bg-slate-100 px-3 py-2 text-sm font-semibold"
              href="/settings/integrations/vercel"
            >
              Vercel
            </Link>
          </li>
        </ul>
      </nav>
      <main className="flex-1 overflow-y-scroll">{children}</main>
    </div>
  );
}
