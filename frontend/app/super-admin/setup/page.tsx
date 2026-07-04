import { redirect } from 'next/navigation'

export default function SuperAdminSetupRedirectPage() {
  redirect('/super-admin/login')
}
