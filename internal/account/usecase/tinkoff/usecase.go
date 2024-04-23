package tinkoff

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
	"go-portfolios-tracker/internal/account"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/internal/logging/slog"
	"go-portfolios-tracker/internal/models"
)

type TestUseCase struct {
	logger *slog.Logger
	config *config.Config
}

type PortfolioUseCase struct {
	logger        logging.Logger
	portfolioRepo account.Repository
}

func NewPortfolioUseCase(logger logging.Logger, portfolioRepo account.Repository) *PortfolioUseCase {
	return &PortfolioUseCase{
		logger:        logger,
		portfolioRepo: portfolioRepo,
	}
}

func (p *PortfolioUseCase) parsePortfolio(tinkoffRepo *investgo.PortfolioResponse) (*models.Portfolio, error) {
	modelPortfolio := &models.Portfolio{}

	err := copier.Copy(modelPortfolio, tinkoffRepo)
	if err != nil {
		return nil, err
	}

	return modelPortfolio, nil
}

func (p *PortfolioUseCase) Init(ctx context.Context, username string, token string) error {
	// TODO: remake type
	err := p.portfolioRepo.Init(ctx, username, token, "tinkoff")
	if err != nil {
		// TODO: error
		return err
	}

	return nil
}

func (p *PortfolioUseCase) Update(ctx context.Context, id int) error {
	/* token, err := p.portfolioRepo.GetToken(ctx, id)
	if err != nil {
		// TODO: error
		return err
	}*/

	account, err := p.portfolioRepo.Get(ctx, id)
	// TODO: tinkoff type
	if err != nil || account.Type != "tinkoff" {
		// TODO
		return err
	}

	// TODO: change config path
	config, err := investgo.LoadConfig("internal/account/config.yaml")
	if err != nil {
		// TODO: return error
		p.logger.Fatal("load config error")
	}
	config.Token = account.Token

	client, err := investgo.NewClient(ctx, config, p.logger)
	if err != nil {
		// TODO: handle error
		p.logger.Fatal("new client error")
	}

	usersService := client.NewUsersServiceClient()
	tinkoffAccounts, err := usersService.GetAccounts()
	if err != nil {
		// TODO: handle error
		p.logger.Fatal("get accounts error")
	}

	operationsService := client.NewOperationsServiceClient()
	account.Portfolios = make([]*models.Portfolio, len(tinkoffAccounts.Accounts))
	for i := 0; i < len(tinkoffAccounts.Accounts); i++ {
		tinkoffRepo, err := operationsService.GetPortfolio(tinkoffAccounts.Accounts[i].Id, pb.PortfolioRequest_RUB)
		if err != nil {
			// TODO: return error
			p.logger.Fatalf("get account error: %v", err)
		}

		account.Portfolios[i], err = p.parsePortfolio(tinkoffRepo)
		if err != nil {
			// TODO: handle error
			p.logger.Fatal(err)
		}
	}

	err = p.portfolioRepo.Update(ctx, account)
	if err != nil {
		// TODO: handle error
		p.logger.Fatal(err)
	}

	return nil
}

func (p *PortfolioUseCase) Get(ctx context.Context, id int) (*models.Account, error) {
	account, err := p.portfolioRepo.Get(ctx, id)
	if err != nil {
		// TODO
		return nil, err
	}

	return account, nil
}

func (p *PortfolioUseCase) Delete(ctx context.Context, id int) error {
	err := p.portfolioRepo.Delete(ctx, id)
	if err != nil {
		// TODO: error
		return err
	}

	return nil
}

func NewTestUseCase() *TestUseCase {
	return &TestUseCase{}
}

func (tuc *TestUseCase) Ha(ctx context.Context, logger logging.Logger) error {
	// TODO: change config path
	config, err := investgo.LoadConfig("internal/account/config.yaml")
	if err != nil {
		// TODO: return error
		logger.Fatal("load config error")
	}

	client, err := investgo.NewClient(ctx, config, logger)
	if err != nil {
		// TODO: handle error
		logger.Fatal("new client error")
	}

	usersService := client.NewUsersServiceClient()
	accounts, err := usersService.GetAccounts()
	if err != nil {
		// TODO: handle error
		logger.Fatal("get accounts error")
	}

	operationsService := client.NewOperationsServiceClient()
	portfolioResp, err := operationsService.GetPortfolio(accounts.Accounts[0].Id, pb.PortfolioRequest_RUB)
	if err != nil {
		// TODO: return error
		logger.Fatalf("get account error: %v", err)
	}

	positions := portfolioResp.Positions
	for i := range positions {
		logger.Info(positions[i])
	}

	return errors.New("HAHAHA")
}
