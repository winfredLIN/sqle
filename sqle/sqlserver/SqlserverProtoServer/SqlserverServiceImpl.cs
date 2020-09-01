﻿using System;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.IO;
using Grpc.Core;
using SqlserverProto;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class SqlServerServiceImpl : SqlserverService.SqlserverServiceBase {
        /*
         * The following compatibility level values can be configured (not all versions supports all of the above listed compatibility level):
         * Product                          Database Engine Version         Compatibility Level Designation         Supported Compatibility Level Values
         * SQL Server 2019 preview                     15                              150                                     150,140,130,120,110,100
         * SQL Server 2017 (14.x)                      14                              140                                     140,130,120,110,100
         * Azure SQL Database logical server           12                              130                                     150,140,130,120,110,100
         * Azure SQL Database Managed Instance         12                              130                                     150,140,130,120,110,100
         * SQL Server 2016 (13.x)                      13                              130                                     130,120,110,100
         * SQL Server 2014 (12.x)                      12                              120                                     120,110,100
         * SQL Server 2012 (11.x)                      11                              110                                     110,100,90
         * SQL Server 2008 R2                          10.5                            100                                     100,90,80
         * SQL Server 2008                             10                              100                                     100,90,80
         * SQL Server 2005 (9.x)                       9                               90                                      90,80
         * SQL Server 2000                             8                               80                                      80
         * 
         * more information: https://docs.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-compatibility-level?view=sql-server-2017
         */
        private const String SQL80 = "80";
        private const String SQL90 = "90";
        private const String SQL100 = "100";
        private const String SQL110 = "110";
        private const String SQL120 = "120";
        private const String SQL130 = "130";

        // sql server parser
        private readonly Dictionary<String, TSqlParser> SqlParsers;

        private TSqlParser GetParser(String version) {
            // set default sql parser version to SQL100
            if (version == "") {
                version = SQL130;
            }

            if (!SqlParsers.ContainsKey(version)) {
                throw new ArgumentException(String.Format("unsupported TSqlParser version:{0}", version));
            }
            return SqlParsers[version];
        }

        // construct function
        public SqlServerServiceImpl() {
            SqlParsers = new Dictionary<String, TSqlParser> {
                {SQL80, new TSql80Parser(false)},
                {SQL90, new TSql90Parser(false)},
                {SQL100, new TSql100Parser(false)},
                {SQL110, new TSql110Parser(false)},
                {SQL120, new TSql120Parser(false)},
                {SQL130, new TSql130Parser(false)},
            };
        }

        internal class SQLVisitor: TSqlFragmentVisitor {
            public List<TSqlStatement> Statements = new List<TSqlStatement>();

            public override void ExplicitVisit(CreateProcedureStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(AlterProcedureStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DropProcedureStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(CreateFunctionStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(AlterFunctionStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DropFunctionStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(UseStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(CreateTableStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(AlterTableStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DropTableStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(CreateDatabaseStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DropDatabaseStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DropSchemaStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(CreateIndexStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DropIndexStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(InsertStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(DeleteStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(UpdateStatement node) {
                Statements.Add(node);
            }

            public override void ExplicitVisit(SelectStatement node) {
                Statements.Add(node);
            }
        }

        // parse sqls
        public IList<TSqlStatement> ParseStatementList(Logger logger, String version, String text) {
            // get parser
            var parser = GetParser(version);

            // parse sqls
            var needRetry = true;
        Try:
            var reader = new StringReader(text);
            IList<ParseError> errorList;
            var statementList = parser.ParseStatementList(reader, out errorList);
            if (errorList.Count > 0) {
                var errMsgList = new List<String>();
                foreach (var parseErr in errorList) {
                    errMsgList.Add(parseErr.Message);
                }

                // remove GO or GO; statement and retry
                if (needRetry) {
                    logger.Info("parse statement error:{0}\nIt will retry after remove GO or GO; statement", String.Join("; ", errMsgList));

                    var sqlLines = text.Split('\n');
                    var newSqlLines = new List<String>();
                    foreach (var sqlLine in sqlLines) {
                        if (sqlLine.Trim().TrimEnd(';').ToUpper() == "GO") {
                            continue;
                        }
                        newSqlLines.Add(sqlLine);
                    }
                    text = String.Join('\n', newSqlLines);

                    needRetry = false;
                    goto Try;
                }

                throw new ArgumentException(String.Format("parse sql `{0}` error: {1}", text, String.Join("; ", errMsgList)));
            }

            return statementList.Statements;
        }

        // parse procedure & function
        public IList<TSqlStatement> ParseWithVisitor(Logger logger, String version, String text) {
            // get parser
            var parser = GetParser(version);

            var reader = new StringReader(text);
            IList<ParseError> errorList;
            TSqlFragment sqlFragment = parser.Parse(reader, out errorList);
            if (errorList.Count > 0) {
                var errMsgList = new List<String>();
                foreach (var parseErr in errorList) {
                    errMsgList.Add(parseErr.Message);
                }
                throw new ArgumentException(String.Format("parse sql `{0}` error: {1}", text, String.Join("; ", errMsgList)));
            }

            SQLVisitor sqlVisitor = new SQLVisitor();
            sqlFragment.Accept(sqlVisitor);

            return sqlVisitor.Statements;
        }

        public IList<TSqlStatement> Parse(Logger logger, String version, String text) {
            try {
                var statements = ParseStatementList(logger, version, text);
                return statements;
            } catch (Exception parseStatementListException) {
                logger.Fatal("ParseStatementList exception error: {0}", parseStatementListException.Message);
                logger.Fatal("ParseStatementList exception stackstrace: {0}", parseStatementListException.StackTrace);

                // ParseStatementList can not parse procedure or function, using Parse try again.
                try {
                    var statements = ParseWithVisitor(logger, version, text);
                    return statements;
                } catch (Exception parseWithVisitorException) {
                    logger.Fatal("ParseWithVisitor exception error: {0}", parseWithVisitorException.Message);
                    logger.Fatal("ParseWithVisitor exception stackstrace: {0}", parseWithVisitorException.StackTrace);

                    throw parseWithVisitorException;
                }
            }
        }

        // Splite sqls
        public override Task<SplitSqlsOutput> GetSplitSqls(SplitSqlsInput request, ServerCallContext context) {
            var output = new SplitSqlsOutput();
            var version = request.Version;
            var sqls = request.Sqls;
            var logger = LogManager.GetCurrentClassLogger();

            try {
                var statements = Parse(logger, version, sqls);

                foreach (var statement in statements) {
                    var sql = "";
                    for (int index = statement.FirstTokenIndex; index <= statement.LastTokenIndex; index++) {
                        sql += statement.ScriptTokenStream[index].Text;
                    }

                    var splitSql = new Sql();
                    splitSql.Sql_ = sql;
                    splitSql.IsDDL = IsDDL(statement);
                    splitSql.IsDML = IsDML(statement);
                    splitSql.IsProcedure = IsProcedure(statement);
                    splitSql.IsFunction = IsFunction(statement);
                    output.SplitSqls.Add(splitSql);
                }
            } catch (Exception e) {
                logger.Fatal("GetSplitSqls exception message:{0}", e.Message);
                logger.Fatal("GetSplitSqls exception stackstrace:{0}", e.StackTrace);
                throw new RpcException(new Status(StatusCode.Internal, e.Message));
            }

            return Task.FromResult(output);
        }

        public bool IsDDL(TSqlStatement statement) {
            if (statement is CreateDatabaseStatement) {
                return true;
            } else if (statement is CreateTableStatement) {
                return true;
            } else if (statement is AlterTableStatement) {
                return true;
            } else if (statement is CreateIndexStatement) {
                return true;
            } else if (statement is DropIndexStatement) {
                return true;
            } else if (statement is DropTableStatement) {
                return true;
            }
            return false;
        }

        public bool IsDML(TSqlStatement statement) {
            if (statement is InsertStatement) {
                return true;
            } else if (statement is UpdateStatement) {
                return true;
            } else if (statement is DeleteStatement) {
                return true;
            }
            return false;
        }

        public bool IsProcedure(TSqlStatement statement) {
            if (statement is CreateProcedureStatement) {
                return true;
            } else if (statement is AlterProcedureStatement) {
                return true;
            } else if (statement is DropProcedureStatement) {
                return true;
            }

            return false;
        }

        public bool IsFunction(TSqlStatement statement) {
            if (statement is CreateFunctionStatement) {
                return true;
            } else if (statement is AlterFunctionStatement) {
                return true;
            } else if (statement is DropFunctionStatement) {
                return true;
            }

            return false;
        }

        // Advise implement
        public override Task<AdviseOutput> Advise(AdviseInput request, ServerCallContext context) {
            var output = new AdviseOutput();
            var version = request.Version;
            var sqls = request.Sqls;
            var ruleNames = request.RuleNames;
            var contextSqls = request.DDLContextSqls;
            var meta = request.SqlserverMeta;
            var contextStart = 0;
            var logger = LogManager.GetCurrentClassLogger();
            logger.Info("advise sqls:{0}\nrules:{1}", String.Join("\n", sqls), String.Join("\n", ruleNames));
            logger.Info("advise host:{0}, port:{1}, user:{2}, current database:{3}", meta.Host, meta.Port, meta.User, meta.CurrentDatabase);

        Try:
            var ruleValidatorContext = new SqlserverContext(meta);
            try {
                for (var index = contextStart; index < contextSqls.Count; index++) {
                    logger.Info("context {0} sqls: {1}", index, String.Join("\n", contextSqls[index].Sqls));
                    foreach (var sql in contextSqls[index].Sqls) {
                        var statements = Parse(logger, version, sql);
                        foreach (var statement in statements) {
                            ruleValidatorContext.UpdateContext(logger, statement);
                        }
                    }
                }
            } catch (Exception e) {
                logger.Fatal("Advise context message:{0}", e.Message);
                logger.Fatal("Advise context exception stacktrace:{0}", e.StackTrace);
                throw new RpcException(new Status(StatusCode.Internal, "parse context error:" + e.Message));
            }
            contextStart++;


            var baseValidatorStatus = AdviseResultContext.BASE_RULE_OK;
            foreach (var sql in sqls) {
                try {
                    var statements = Parse(logger, version, sql);
                    foreach (var statement in statements) {
                        var baseRuleValidator = new BaseRuleValidator();
                        baseRuleValidator.Check(ruleValidatorContext, statement);
                        if (ruleValidatorContext.AdviseResultContext.GetBaseRuleStatus() == AdviseResultContext.BASE_RULE_FAILED) {
                            baseValidatorStatus = AdviseResultContext.BASE_RULE_FAILED;
                        }

                        foreach (var ruleName in ruleNames) {
                            if (!DefaultRules.RuleValidators.ContainsKey(ruleName)) {
                                continue;
                            }
                            DefaultRules.RuleValidators[ruleName].Check(ruleValidatorContext, statement);
                        }

                        ruleValidatorContext.UpdateContext(logger, statement);
                        var result = ruleValidatorContext.AdviseResultContext.GetAdviseResult();
                        result.IsDDL = IsDDL(statement);
                        result.IsDML = IsDML(statement);
                        result.IsProcedure = IsProcedure(statement);
                        result.IsFunction = IsFunction(statement);
                        output.Results.Add(result);
                    }
                    ruleValidatorContext.AdviseResultContext.ResetAdviseResult();
                } catch (Exception e) {
                    logger.Fatal("Advise exception stacktrace:{0}", e.StackTrace);
                    logger.Fatal("Advise exception message:{0}", e.Message);
                    throw new RpcException(new Status(StatusCode.Internal, e.Message));
                }
            }

            if (baseValidatorStatus == AdviseResultContext.BASE_RULE_FAILED && contextStart < contextSqls.Count) {
                output.Results.Clear();
                goto Try;
            }

            if (baseValidatorStatus == AdviseResultContext.BASE_RULE_FAILED) {
                output.BaseValidatorFailed = true;
            }

            return Task.FromResult(output);
        }

        // GetRollbackSqls implement
        public override Task<GetRollbackSqlsOutput> GetRollbackSqls(GetRollbackSqlsInput request, ServerCallContext context) {
            var output = new GetRollbackSqlsOutput();
            var version = request.Version;
            var sqls = request.Sqls;
            var meta = request.SqlserverMeta;
            var rollbackSqlContext = new SqlserverContext(meta, request.RollbackConfig);
            var logger = LogManager.GetCurrentClassLogger();
            logger.Info("getrollback host:{0}, port:{1}, user:{2}, current database:{3}", meta.Host, meta.Port, meta.User, meta.CurrentDatabase);

            foreach (var sql in sqls) {
                try {
                    var statements = Parse(logger, version, sql);
                    logger.Info("sql:{0}\n", sql);
                    foreach (var statement in statements) {
                        if (IsProcedure(statement)) {
                            var procedureRollbackSql = new Sql {
                                IsProcedure = true,
                            };
                            output.RollbackSqls.Add(procedureRollbackSql);
                        } else if (IsFunction(statement)) {
                            var functionRollbackSql = new Sql {
                                IsFunction = true,
                            };
                            output.RollbackSqls.Add(functionRollbackSql);
                        } else {
                            output.RollbackSqls.Add(new RollbackSql().GetRollbackSql(rollbackSqlContext, statement));
                            rollbackSqlContext.UpdateContext(logger, statement);
                        }
                    }
                } catch (Exception e) {
                    logger.Fatal("GetRollbackSqls exception message:{0}", e.Message);
                    logger.Fatal("GetRollbackSqls exception stackstrace:{0}", e.StackTrace);
                    throw new RpcException(new Status(StatusCode.Internal, e.Message));
                }

            }

            return Task.FromResult(output);
        }

        // GetProcedureFunctionBackupSql implement
        public override Task<GetProcedureFunctionBackupSqlOutput> GetProcedureFunctionBackupSql(GetProcedureFunctionBackupSqlInput request, ServerCallContext context) {
            var output = new GetProcedureFunctionBackupSqlOutput();
            var logger = LogManager.GetCurrentClassLogger();
            var sqlserverMeta = request.SqlserverMeta;
            ProcedureFunctionBackupSqlGenerator procedureFunctionBackupSqlGenerator = new ProcedureFunctionBackupSqlGenerator(request.SqlserverMeta);

            try {
                var statements = Parse(logger, "", request.Sql);
                foreach (var statement in statements) {
                    List<string> backupSqls = new List<string>();
                    // drop procedure
                    if (statement is DropProcedureStatement) {
                        var dropProcedureStatement = statement as DropProcedureStatement;
                        backupSqls = procedureFunctionBackupSqlGenerator.GetbackupSqlsForObejcts(logger, dropProcedureStatement.Objects);
                    }

                    // alter procedure
                    if (statement is AlterProcedureStatement) {
                        var alterProcedureStatement = statement as AlterProcedureStatement;
                        var objects = new List<SchemaObjectName>();
                        objects.Add(alterProcedureStatement.ProcedureReference.Name);
                        backupSqls = procedureFunctionBackupSqlGenerator.GetbackupSqlsForObejcts(logger, objects);
                    }

                    // drop function
                    if (statement is DropFunctionStatement) {
                        var dropFunctionStatement = statement as DropFunctionStatement;
                        backupSqls = procedureFunctionBackupSqlGenerator.GetbackupSqlsForObejcts(logger, dropFunctionStatement.Objects);
                    }

                    // alter function
                    if (statement is AlterFunctionStatement) {
                        var alterFunctionStatement = statement as AlterFunctionStatement;
                        var obejcts = new List<SchemaObjectName>();
                        obejcts.Add(alterFunctionStatement.Name);
                        backupSqls = procedureFunctionBackupSqlGenerator.GetbackupSqlsForObejcts(logger, obejcts);
                    }

                    foreach (var backupSql in backupSqls) {
                        if (!output.BackupSqls.Contains(backupSql)) {
                            output.BackupSqls.Add(backupSql);
                        }
                    }
                }
            } catch(Exception e) {
                logger.Fatal("GetProcedureFunctionBackupSql exception message:{0}", e.Message);
                logger.Fatal("GetProcedureFunctionBackupSql exception stackstrace:{0}", e.StackTrace);

                throw new RpcException(new Status(StatusCode.Internal, e.Message));
            }
            return Task.FromResult(output);
        }
    }
}