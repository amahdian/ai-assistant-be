package storage

import "context"

type Storage interface {
	User(ctx context.Context) UserStorage
}

type Session interface {
	// Begin starts a transactional session.
	//
	// It's the user's responsibility to manage the session,
	// Either Rollback or Commit MUST be called to pair with Begin to avoid transaction leak.
	Begin() (Session, error)
	// Rollback aborts the changes made by the transactional session.
	Rollback() error
	// Commit commits the changes made by the transactional session.
	Commit() error
	// Close the session
	Close() error
}

var sessionKeyInCtx = "env:ctx:storage_session"

func WithContext(ctx context.Context, ses Session) context.Context {
	return context.WithValue(ctx, sessionKeyInCtx, ses)
}

func FromContext(ctx context.Context) Session {
	v := ctx.Value(sessionKeyInCtx)
	if v == nil {
		return nil
	}
	return v.(Session)
}
