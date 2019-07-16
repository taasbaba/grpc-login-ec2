package v1

/*
 * 在這個檔案裡實作 poker-service.proto 內定義的 inter
*/
import (
	"context"
	"database/sql"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	auth "grpc-login-server/server/internal/api/v1"
	"grpc-login-server/server/internal/logger"
	sqlWrapper "grpc-login-server/server/internal/sqlwrapper"
	"strings"
)
type theUser struct {
	Id string
	Pd string
}

type Server struct{
	User []theUser
	db *sqlWrapper.DB
}

func NewAuthServer(db *sql.DB) auth.AuthServer {
	sw := sqlWrapper.WrapperDB(db, true, 1)
	return &Server{db:sw}
}

func (s *Server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT `password` FROM user WHERE `email` = ?",
		req.Username)
	if err != nil {
		logger.Log.Info("Login:QueryContext:", zap.String("err", err.Error()))
		return &auth.LoginResponse{Token: "fail"}, nil
	}
	defer rows.Close()

	if !rows.Next() { // can't find username
		if err := rows.Err(); err != nil {
			logger.Log.Info("Login:!rows.Next():", zap.String("err", err.Error()))
			return &auth.LoginResponse{Token: "fail"}, nil
		}
	} else {
		var password string
		if err := rows.Scan(&password); err != nil {
			logger.Log.Info("Login:err := rows.Scan(&password):", zap.String("err", err.Error()))
			return &auth.LoginResponse{Token: "fail"}, nil
		}

		logger.Log.Info("Login:", zap.String("hashpassword", password))

		if comparePasswords(password, []byte(req.Password)) {
			return &auth.LoginResponse{Token: "ok"}, nil
		} else {
			return &auth.LoginResponse{Token: "fail"}, nil
		}
	}

	return &auth.LoginResponse{Token: "fail"}, nil

}

func (s *Server) Registration(ctx context.Context, req *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	// get SQL connection from pool
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := s.db.QueryContext(ctx, "SELECT `email` FROM user WHERE `email` = ?",
		req.Username)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from user "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() { // username not use
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from user"+err.Error())
		}
	} else {
		logger.Log.Info("Registration::fail repeat:", zap.String("id:", req.Username), zap.String("pd:", req.Password))
		return &auth.RegistrationResponse{Token: "Repeat"}, nil
	}

	hash := hashAndSalt([]byte(req.Password))
	logger.Log.Info("Registration::success:", zap.String("id:", req.Username), zap.String("pd:", req.Password))
	s.User = append(s.User, theUser{req.Username, req.Password})

	userId := strings.Replace(uuid.NewV4().String(), "-", "", -1 )
	// insert username to user
	res, err := s.db.ExecContext(ctx, "INSERT INTO user(`user_id`, `nickname`, `email`, `password`) VALUES(?, ?, ?, ?)",
		userId, "", req.Username, hash)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into user: "+err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created user: "+err.Error())
	}
	logger.Log.Info("Registration:", zap.Int64("id:", id))
	return &auth.RegistrationResponse{Token: "success"}, nil

}


func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		logger.Log.Info("hashAndSalt:", zap.String("err", err.Error()))
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		logger.Log.Info("comparePasswords:", zap.String("err", err.Error()))
		return false
	}

	return true
}
