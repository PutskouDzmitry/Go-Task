package wallet

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"task/internal/entity"
	"task/pkg/api/filter"
	"task/pkg/database/postgresql"
)

func (w *Wallet) CreateWallet(ctx context.Context, coin string) error {
	var startCoin = "000000000000000000"
	sql, args, err := w.queryBuilder.Insert("wallet").Columns("coin", "money").Values(coin, startCoin).ToSql()

	logging := w.log.With(
		zap.String("sql", sql),
		zap.Any("args", args),
	)

	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		logging.Error(err.Error())
		return err
	}

	exec, err := w.db.Exec(ctx, sql, args...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	if exec.RowsAffected() == 0 || !exec.Insert() {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}
	return nil
}

func (w *Wallet) GetWallet(ctx context.Context, fo *filter.Options) (*entity.WalletDB, error) {
	walletBalance := &entity.WalletDB{}
	columns := walletBalance.Columns("wallet")

	q := w.queryBuilder.
		Select(columns...).
		From("wallet")

	q = postgresql.NewFilterOptions(fo).Enrich(q, "wallet")

	sql, args, err := q.ToSql()
	logger := w.log.With(
		zap.String("sql", sql),
		zap.Any("args", args),
	)

	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		logger.Error("failed to get wallet" + err.Error())
		return nil, err
	}

	if err = w.db.QueryRow(ctx, sql, args...).Scan(walletBalance.Dest()...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		err = postgresql.ErrScan(err)
		logger.Error("failed to get wallet" + err.Error())
		return nil, err
	}
	return walletBalance, nil
}

func (w *Wallet) IncOrDecrBalance(ctx context.Context, m map[string]interface{}, id string, mWallet map[string]interface{}) error {
	tx, err := w.db.Begin(ctx)
	if err != nil {
		err = postgresql.ErrCreateTx(err)
		w.log.Error("failed to create transaction" + err.Error())
		return err
	}

	defer func() {
		if err != nil {
			w.log.Error(err.Error())
			err = tx.Rollback(ctx)
			if err != nil {
				err = postgresql.ErrRollback(err)
				w.log.Error("failed to rollback transaction" + err.Error())
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				err = postgresql.ErrCommit(err)
				w.log.Error("failed to commit transaction" + err.Error())
				return
			}
		}
	}()

	sql, args, err := w.queryBuilder.Insert("transactional").SetMap(m).ToSql()
	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		w.log.Error("failed to create online school" + err.Error())
		return err
	}

	logging := w.log.With(
		zap.String("sql", sql),
		zap.Any("args", args),
	)

	// Update wallet
	sqlWallet, argsWallet, err := w.queryBuilder.
		Update("wallet").
		SetMap(mWallet).
		Where(sq.Eq{"id": id}).
		ToSql()

	logger := w.log.With(
		zap.String("sql", sqlWallet),
		zap.Any("args", argsWallet),
	)

	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		logger.Error("failed to update wallet" + err.Error())
		return err
	}

	// trans
	exec, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	if exec.RowsAffected() == 0 || !exec.Insert() {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	// wallet
	exec, err = tx.Exec(ctx, sqlWallet, argsWallet...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logger.Error("failed to update wallet" + err.Error())
		return err
	}

	if exec.RowsAffected() != 1 || !exec.Update() {
		err = postgresql.ErrDoQuery(err)
		logger.Error("missing update wallet. Updating 0 rows" + err.Error())
		return err
	}

	return nil
}

func (w *Wallet) MoneyTransactional(ctx context.Context, mBalanceFrom, mTransFrom, mBalanceTo, mTransTo map[string]interface{}) error {
	tx, err := w.db.Begin(ctx)
	if err != nil {
		err = postgresql.ErrCreateTx(err)
		w.log.Error("failed to create transaction" + err.Error())
		return err
	}

	defer func() {
		if err != nil {
			w.log.Error(err.Error())
			err = tx.Rollback(ctx)
			if err != nil {
				err = postgresql.ErrRollback(err)
				w.log.Error("failed to rollback transaction" + err.Error())
				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				err = postgresql.ErrCommit(err)
				w.log.Error("failed to commit transaction" + err.Error())
				return
			}
		}
	}()

	// trans from
	sqlTransFrom, argsTransFrom, err := w.queryBuilder.Insert("transactional").SetMap(mTransFrom).ToSql()
	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		w.log.Error("failed to create online school" + err.Error())
		return err
	}

	logging := w.log.With(
		zap.String("sql", sqlTransFrom),
		zap.Any("args", argsTransFrom),
	)

	// Update wallet from
	sqlWalletFrom, argsWalletFrom, err := w.queryBuilder.
		Update("wallet").
		SetMap(mBalanceFrom).
		Where(sq.Eq{"id": mBalanceFrom["id"]}).
		ToSql()

	logger := w.log.With(
		zap.String("sql", sqlWalletFrom),
		zap.Any("args", argsWalletFrom),
	)

	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		logger.Error("failed to update wallet" + err.Error())
		return err
	}

	// trans to
	sqlTransTo, argsTransTo, err := w.queryBuilder.Insert("transactional").SetMap(mTransTo).ToSql()
	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		w.log.Error("failed to create online school" + err.Error())
		return err
	}

	logging = w.log.With(
		zap.String("sql", sqlTransTo),
		zap.Any("args", argsTransTo),
	)

	// Update wallet from
	sqlWalletTo, argsWalletTo, err := w.queryBuilder.
		Update("wallet").
		SetMap(mBalanceTo).
		Where(sq.Eq{"id": mBalanceTo["id"]}).
		ToSql()

	logger = w.log.With(
		zap.String("sql", sqlWalletTo),
		zap.Any("args", argsWalletTo),
	)

	if err != nil {
		err = postgresql.ErrCreateQuery(err)
		logger.Error("failed to update wallet" + err.Error())
		return err
	}

	// trans from
	exec, err := tx.Exec(ctx, sqlTransFrom, argsTransFrom...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	if exec.RowsAffected() == 0 || !exec.Insert() {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	// wallet from
	exec, err = tx.Exec(ctx, sqlWalletFrom, argsWalletFrom...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logger.Error("failed to update wallet" + err.Error())
		return err
	}

	if exec.RowsAffected() != 1 || !exec.Update() {
		err = postgresql.ErrDoQuery(err)
		logger.Error("missing update wallet. Updating 0 rows" + err.Error())
		return err
	}

	// trans to
	exec, err = tx.Exec(ctx, sqlTransTo, argsTransTo...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	if exec.RowsAffected() == 0 || !exec.Insert() {
		err = postgresql.ErrDoQuery(err)
		logging.Error(err.Error())
		return err
	}

	// wallet from
	exec, err = tx.Exec(ctx, sqlWalletTo, argsWalletTo...)
	if err != nil {
		err = postgresql.ErrDoQuery(err)
		logger.Error("failed to update wallet" + err.Error())
		return err
	}

	if exec.RowsAffected() != 1 || !exec.Update() {
		err = postgresql.ErrDoQuery(err)
		logger.Error("missing update wallet. Updating 0 rows" + err.Error())
		return err
	}

	return nil
}
