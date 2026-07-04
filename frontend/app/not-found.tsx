import Link from 'next/link'

export default function NotFound() {
  return (
    <main
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '24px',
        background:
          'radial-gradient(circle at top, rgba(55,138,221,0.18), transparent 28%), linear-gradient(180deg, #07121F 0%, #0B1629 100%)',
      }}
    >
      <div
        style={{
          width: '100%',
          maxWidth: '560px',
          borderRadius: '24px',
          border: '1px solid rgba(120,140,200,0.18)',
          background: 'rgba(17,30,53,0.82)',
          boxShadow: '0 30px 70px rgba(0, 0, 0, 0.28)',
          padding: '36px 30px',
          color: '#E8EAF2',
          textAlign: 'center',
          backdropFilter: 'blur(10px)',
        }}
      >
        <div
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: '72px',
            height: '72px',
            marginBottom: '18px',
            borderRadius: '999px',
            background: 'rgba(239,159,39,0.14)',
            color: '#FAC775',
            fontSize: '28px',
            fontWeight: 800,
          }}
        >
          404
        </div>

        <h1 style={{ margin: '0 0 12px', fontSize: '30px', fontWeight: 800 }}>
          페이지를 찾을 수 없습니다
        </h1>
        <p
          style={{
            margin: '0 0 24px',
            fontSize: '15px',
            lineHeight: 1.7,
            color: 'rgba(200,210,235,0.76)',
          }}
        >
          요청하신 주소가 변경되었거나 존재하지 않습니다.
          <br />
          아래 버튼을 눌러 메인 페이지로 돌아가세요.
        </p>

        <Link
          href="/"
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            minWidth: '180px',
            padding: '13px 20px',
            borderRadius: '999px',
            background: 'linear-gradient(135deg, #378ADD 0%, #7F77DD 100%)',
            color: '#FFFFFF',
            textDecoration: 'none',
            fontSize: '14px',
            fontWeight: 800,
            boxShadow: '0 12px 24px rgba(55,138,221,0.22)',
          }}
        >
          메인 페이지로 이동
        </Link>
      </div>
    </main>
  )
}
